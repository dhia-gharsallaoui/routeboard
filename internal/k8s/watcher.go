package k8s

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwinformers "sigs.k8s.io/gateway-api/pkg/client/informers/externalversions"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/store"
)

type Watcher struct {
	cfg     *config.Config
	clients *Clients
	store   *store.Store
	queue   workqueue.TypedRateLimitingInterface[string]

	k8sFactory k8sinformers.SharedInformerFactory
	gwFactory  gwinformers.SharedInformerFactory

	ingressInformer   cache.SharedIndexInformer
	httprouteInformer cache.SharedIndexInformer
}

func NewWatcher(cfg *config.Config, clients *Clients, store *store.Store) *Watcher {
	queue := workqueue.NewTypedRateLimitingQueue(
		workqueue.DefaultTypedControllerRateLimiter[string](),
	)

	w := &Watcher{
		cfg:     cfg,
		clients: clients,
		store:   store,
		queue:   queue,
	}

	if cfg.WatchIngress {
		w.k8sFactory = k8sinformers.NewSharedInformerFactory(clients.Kubernetes, cfg.ResyncInterval)
		w.ingressInformer = w.k8sFactory.Networking().V1().Ingresses().Informer()
		if _, err := w.ingressInformer.AddEventHandler(w.eventHandler(string(sourceIngress))); err != nil {
			slog.Error("failed to add ingress event handler", "error", err)
		}
	}

	if cfg.WatchHTTPRoute {
		w.gwFactory = gwinformers.NewSharedInformerFactory(clients.GatewayAPI, cfg.ResyncInterval)
		w.httprouteInformer = w.gwFactory.Gateway().V1().HTTPRoutes().Informer()
		if _, err := w.httprouteInformer.AddEventHandler(w.eventHandler(string(sourceHTTPRoute))); err != nil {
			slog.Error("failed to add httproute event handler", "error", err)
		}
	}

	return w
}

const (
	sourceIngress   = "Ingress"
	sourceHTTPRoute = "HTTPRoute"
)

func (w *Watcher) eventHandler(source string) cache.ResourceEventHandlerFuncs {
	enqueue := func(obj interface{}) {
		key, err := cache.MetaNamespaceKeyFunc(obj)
		if err != nil {
			runtime.HandleError(err)
			return
		}
		w.queue.Add(source + ":" + key)
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc:    enqueue,
		UpdateFunc: func(_, newObj interface{}) { enqueue(newObj) },
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				runtime.HandleError(err)
				return
			}
			w.store.Delete(source + ":" + key)
		},
	}
}

func (w *Watcher) Run(ctx context.Context) error {
	defer runtime.HandleCrash()
	defer w.queue.ShutDown()

	slog.Info("starting watcher")

	if w.k8sFactory != nil {
		w.k8sFactory.Start(ctx.Done())
	}
	if w.gwFactory != nil {
		w.gwFactory.Start(ctx.Done())
	}

	syncs := make([]cache.InformerSynced, 0, 2)
	if w.ingressInformer != nil {
		syncs = append(syncs, w.ingressInformer.HasSynced)
	}
	if w.httprouteInformer != nil {
		syncs = append(syncs, w.httprouteInformer.HasSynced)
	}
	if !cache.WaitForCacheSync(ctx.Done(), syncs...) {
		return fmt.Errorf("timed out waiting for informer caches to sync")
	}

	slog.Info("informer caches synced, starting workers")

	for i := 0; i < 2; i++ {
		go wait.UntilWithContext(ctx, w.runWorker, time.Second)
	}

	<-ctx.Done()
	slog.Info("watcher stopped")
	return nil
}

func (w *Watcher) runWorker(ctx context.Context) {
	for w.processNextItem() {
	}
}

func (w *Watcher) processNextItem() bool {
	key, quit := w.queue.Get()
	if quit {
		return false
	}
	defer w.queue.Done(key)

	err := w.sync(key)
	if err == nil {
		w.queue.Forget(key)
		return true
	}

	if w.queue.NumRequeues(key) < 5 {
		slog.Warn("error syncing, retrying", "key", key, "error", err)
		w.queue.AddRateLimited(key)
	} else {
		slog.Error("dropping key after max retries", "key", key, "error", err)
		w.queue.Forget(key)
	}
	return true
}

func (w *Watcher) sync(compositeKey string) error {
	source, key, ok := strings.Cut(compositeKey, ":")
	if !ok {
		return fmt.Errorf("invalid composite key: %s", compositeKey)
	}

	namespace, _, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return fmt.Errorf("splitting key %s: %w", key, err)
	}

	if !w.shouldInclude(namespace) {
		return nil
	}

	switch source {
	case sourceIngress:
		return w.syncIngress(key)
	case sourceHTTPRoute:
		return w.syncHTTPRoute(key)
	default:
		return fmt.Errorf("unknown source: %s", source)
	}
}

func (w *Watcher) syncIngress(key string) error {
	obj, exists, err := w.ingressInformer.GetStore().GetByKey(key)
	if err != nil {
		return fmt.Errorf("getting ingress %s: %w", key, err)
	}
	if !exists {
		w.store.Delete("Ingress:" + key)
		return nil
	}

	ingress, ok := obj.(*networkingv1.Ingress)
	if !ok {
		return fmt.Errorf("unexpected type for ingress %s", key)
	}

	route := extractIngressRoute(ingress)
	w.store.Set(route)
	return nil
}

func (w *Watcher) syncHTTPRoute(key string) error {
	obj, exists, err := w.httprouteInformer.GetStore().GetByKey(key)
	if err != nil {
		return fmt.Errorf("getting httproute %s: %w", key, err)
	}
	if !exists {
		w.store.Delete("HTTPRoute:" + key)
		return nil
	}

	hr, ok := obj.(*gwv1.HTTPRoute)
	if !ok {
		return fmt.Errorf("unexpected type for httproute %s", key)
	}

	route := extractHTTPRouteRoute(hr)
	w.store.Set(route)
	return nil
}

func (w *Watcher) shouldInclude(namespace string) bool {
	if len(w.cfg.NamespaceAllowlist) > 0 {
		for _, ns := range w.cfg.NamespaceAllowlist {
			if ns == namespace {
				return true
			}
		}
		return false
	}
	for _, ns := range w.cfg.NamespaceDenylist {
		if ns == namespace {
			return false
		}
	}
	return true
}
