package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	gwclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type Clients struct {
	Kubernetes kubernetes.Interface
	GatewayAPI gwclient.Interface
}

func NewClients(kubeconfig string) (*Clients, error) {
	cfg, err := buildConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("building kubeconfig: %w", err)
	}

	k8s, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	gw, err := gwclient.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating gateway API client: %w", err)
	}

	return &Clients{Kubernetes: k8s, GatewayAPI: gw}, nil
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	cfg, err := rest.InClusterConfig()
	if err != nil {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			loadingRules, &clientcmd.ConfigOverrides{},
		).ClientConfig()
	}
	return cfg, nil
}
