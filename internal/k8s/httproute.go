package k8s

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/dhia/routeboard/internal/model"
)

// gatewayLookupFunc resolves a Gateway by namespace/name, typically backed by
// the watcher's Gateway informer cache. It returns nil when the Gateway is
// not available (not cached yet, missing RBAC, or not watched at all).
type gatewayLookupFunc func(namespace, name string) *gwv1.Gateway

func extractHTTPRouteRoute(hr *gwv1.HTTPRoute, lookupGateway gatewayLookupFunc) *model.Route {
	r := &model.Route{
		ID:          fmt.Sprintf("HTTPRoute:%s/%s", hr.Namespace, hr.Name),
		Name:        hr.Name,
		Namespace:   hr.Namespace,
		Source:      model.SourceHTTPRoute,
		Group:       hr.Namespace,
		Labels:      hr.Labels,
		Annotations: hr.Annotations,
		CreatedAt:   hr.CreationTimestamp.Time,
		UpdatedAt:   time.Now(),
	}

	for _, h := range hr.Spec.Hostnames {
		r.Hosts = append(r.Hosts, string(h))
	}

	for _, parent := range hr.Spec.ParentRefs {
		if parentRefTLS(hr, parent, lookupGateway) {
			r.TLS = true
		}
	}

	for _, rule := range hr.Spec.Rules {
		for _, match := range rule.Matches {
			if match.Path != nil && match.Path.Value != nil {
				r.Paths = append(r.Paths, *match.Path.Value)
			}
		}
		for _, backend := range rule.BackendRefs {
			if r.ServiceName == "" {
				r.ServiceName = string(backend.Name)
				if backend.Port != nil {
					r.ServicePort = fmt.Sprintf("%d", *backend.Port)
				}
			}
		}
	}

	r.URL = computeURL(r)
	r.Title = titleize(r.Name)
	r.Icon = model.DetectIcon(r.ServiceName, r.Name)
	model.ApplyAnnotations(r, hr.Annotations)

	return r
}

// parentRefTLS reports whether a parentRef implies the route is served over
// TLS. It resolves the referenced Gateway's listeners when available and
// falls back to the parentRef heuristic otherwise.
func parentRefTLS(hr *gwv1.HTTPRoute, parent gwv1.ParentReference, lookupGateway gatewayLookupFunc) bool {
	if !refersToGateway(parent) {
		return false
	}

	if lookupGateway != nil {
		ns := hr.Namespace
		if parent.Namespace != nil && *parent.Namespace != "" {
			ns = string(*parent.Namespace)
		}
		if gw := lookupGateway(ns, string(parent.Name)); gw != nil {
			return gatewayListenerTLS(gw, parent, hr.Spec.Hostnames)
		}
		slog.Debug("gateway not found in cache, falling back to parentRef heuristic",
			"gateway", ns+"/"+string(parent.Name),
			"httproute", hr.Namespace+"/"+hr.Name)
	}

	return isTLSParentRef(parent)
}

// refersToGateway reports whether a parentRef targets a Gateway API Gateway.
// Empty group/kind default to Gateway per the Gateway API spec.
func refersToGateway(parent gwv1.ParentReference) bool {
	if parent.Group != nil && *parent.Group != "" && string(*parent.Group) != gwv1.GroupName {
		return false
	}
	if parent.Kind != nil && *parent.Kind != "" && *parent.Kind != "Gateway" {
		return false
	}
	return true
}

// gatewayListenerTLS reports whether the parentRef selects at least one
// TLS-terminating listener on the Gateway. Listeners are selected by
// sectionName when set, otherwise by port when set, otherwise all listeners
// are considered (preferring HTTPS when listeners are mixed).
func gatewayListenerTLS(gw *gwv1.Gateway, parent gwv1.ParentReference, routeHostnames []gwv1.Hostname) bool {
	for _, l := range gw.Spec.Listeners {
		if parent.SectionName != nil && l.Name != *parent.SectionName {
			continue
		}
		if parent.SectionName == nil && parent.Port != nil && l.Port != *parent.Port {
			continue
		}
		if l.Protocol != gwv1.HTTPSProtocolType && l.Protocol != gwv1.TLSProtocolType {
			continue
		}
		if !listenerHostnameCompatible(l.Hostname, routeHostnames) {
			continue
		}
		return true
	}
	return false
}

// listenerHostnameCompatible reports whether a listener hostname can serve
// any of the route's hostnames. A nil/empty listener hostname or a route
// without hostnames matches everything. Wildcards match a leading label
// suffix (e.g. "*.example.com" serves "app.example.com").
func listenerHostnameCompatible(listenerHostname *gwv1.Hostname, routeHostnames []gwv1.Hostname) bool {
	if listenerHostname == nil || *listenerHostname == "" || len(routeHostnames) == 0 {
		return true
	}
	lh := string(*listenerHostname)
	for _, rh := range routeHostnames {
		h := string(rh)
		if h == lh {
			return true
		}
		if strings.HasPrefix(lh, "*.") && strings.HasSuffix(h, lh[1:]) {
			return true
		}
		if strings.HasPrefix(h, "*.") && strings.HasSuffix(lh, h[1:]) {
			return true
		}
	}
	return false
}

// isTLSParentRef is the fallback heuristic used when the parent Gateway is
// not available in the cache (not synced yet, missing RBAC on existing
// installs): HTTPRoutes carry no TLS config themselves, so this relies on the
// listener port and common listener naming conventions (e.g. "https", "tls",
// Traefik's "websecure").
func isTLSParentRef(parent gwv1.ParentReference) bool {
	if parent.Port != nil && *parent.Port == 443 {
		return true
	}
	if parent.SectionName == nil {
		return false
	}
	sn := strings.ToLower(string(*parent.SectionName))
	if strings.Contains(sn, "https") || strings.Contains(sn, "tls") {
		return true
	}
	return strings.Contains(sn, "secure") && !strings.Contains(sn, "insecure")
}
