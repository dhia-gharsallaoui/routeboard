package k8s

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func sectionName(s string) *gwv1.SectionName {
	sn := gwv1.SectionName(s)
	return &sn
}

func portNumber(p int32) *gwv1.PortNumber {
	pn := gwv1.PortNumber(p)
	return &pn
}

func hostname(h string) *gwv1.Hostname {
	hn := gwv1.Hostname(h)
	return &hn
}

func namespaceRef(ns string) *gwv1.Namespace {
	n := gwv1.Namespace(ns)
	return &n
}

func gateway(namespace, name string, listeners ...gwv1.Listener) *gwv1.Gateway {
	return &gwv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec:       gwv1.GatewaySpec{Listeners: listeners},
	}
}

// lookupFrom returns a gatewayLookupFunc backed by a fixed set of Gateways,
// mirroring how the watcher resolves Gateways from its informer cache.
func lookupFrom(gateways ...*gwv1.Gateway) gatewayLookupFunc {
	return func(namespace, name string) *gwv1.Gateway {
		for _, gw := range gateways {
			if gw.Namespace == namespace && gw.Name == name {
				return gw
			}
		}
		return nil
	}
}

func newHTTPRoute(parentRefs []gwv1.ParentReference) *gwv1.HTTPRoute {
	return &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "default"},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: parentRefs},
			Hostnames:       []gwv1.Hostname{"app.example.com"},
		},
	}
}

func TestExtractHTTPRouteRouteGatewayListenerTLS(t *testing.T) {
	httpListener := gwv1.Listener{Name: "web", Port: 80, Protocol: gwv1.HTTPProtocolType}
	httpsListener := gwv1.Listener{Name: "websecure", Port: 443, Protocol: gwv1.HTTPSProtocolType}
	tlsListener := gwv1.Listener{Name: "passthrough", Port: 8443, Protocol: gwv1.TLSProtocolType}

	tests := []struct {
		name       string
		parentRefs []gwv1.ParentReference
		gateways   []*gwv1.Gateway
		wantTLS    bool
		wantURL    string
	}{
		{
			name: "sectionName matching HTTPS listener",
			parentRefs: []gwv1.ParentReference{
				{Name: "gw", SectionName: sectionName("websecure")},
			},
			gateways: []*gwv1.Gateway{gateway("default", "gw", httpListener, httpsListener)},
			wantTLS:  true,
			wantURL:  "https://app.example.com",
		},
		{
			name: "sectionName matching HTTP listener stays http even if name sounds secure",
			parentRefs: []gwv1.ParentReference{
				{Name: "gw", SectionName: sectionName("websecure")},
			},
			gateways: []*gwv1.Gateway{gateway("default", "gw",
				gwv1.Listener{Name: "websecure", Port: 80, Protocol: gwv1.HTTPProtocolType})},
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
		{
			name: "port matching HTTPS listener",
			parentRefs: []gwv1.ParentReference{
				{Name: "gw", Port: portNumber(443)},
			},
			gateways: []*gwv1.Gateway{gateway("default", "gw", httpListener, httpsListener)},
			wantTLS:  true,
			wantURL:  "https://app.example.com",
		},
		{
			name: "port matching HTTP listener stays http",
			parentRefs: []gwv1.ParentReference{
				{Name: "gw", Port: portNumber(80)},
			},
			gateways: []*gwv1.Gateway{gateway("default", "gw", httpListener, httpsListener)},
			wantTLS:  false,
			wantURL:  "http://app.example.com",
		},
		{
			name:       "no sectionName or port with mixed listeners prefers https",
			parentRefs: []gwv1.ParentReference{{Name: "gw"}},
			gateways:   []*gwv1.Gateway{gateway("default", "gw", httpListener, httpsListener)},
			wantTLS:    true,
			wantURL:    "https://app.example.com",
		},
		{
			name:       "no sectionName or port with only HTTP listener stays http",
			parentRefs: []gwv1.ParentReference{{Name: "gw"}},
			gateways:   []*gwv1.Gateway{gateway("default", "gw", httpListener)},
			wantTLS:    false,
			wantURL:    "http://app.example.com",
		},
		{
			name:       "TLS protocol listener counts as tls",
			parentRefs: []gwv1.ParentReference{{Name: "gw"}},
			gateways:   []*gwv1.Gateway{gateway("default", "gw", tlsListener)},
			wantTLS:    true,
			wantURL:    "https://app.example.com",
		},
		{
			name: "cross-namespace parentRef resolves gateway in that namespace",
			parentRefs: []gwv1.ParentReference{
				{Name: "gw", Namespace: namespaceRef("infra")},
			},
			gateways: []*gwv1.Gateway{gateway("infra", "gw", httpsListener)},
			wantTLS:  true,
			wantURL:  "https://app.example.com",
		},
		{
			name:       "HTTPS listener with matching wildcard hostname",
			parentRefs: []gwv1.ParentReference{{Name: "gw"}},
			gateways: []*gwv1.Gateway{gateway("default", "gw",
				gwv1.Listener{Name: "https", Port: 443, Protocol: gwv1.HTTPSProtocolType, Hostname: hostname("*.example.com")})},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name:       "HTTPS listener with incompatible hostname stays http",
			parentRefs: []gwv1.ParentReference{{Name: "gw"}},
			gateways: []*gwv1.Gateway{gateway("default", "gw",
				gwv1.Listener{Name: "https", Port: 443, Protocol: gwv1.HTTPSProtocolType, Hostname: hostname("*.other.com")})},
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := newHTTPRoute(tt.parentRefs)

			r := extractHTTPRouteRoute(hr, lookupFrom(tt.gateways...))

			if r.TLS != tt.wantTLS {
				t.Errorf("TLS = %v, want %v", r.TLS, tt.wantTLS)
			}
			if r.URL != tt.wantURL {
				t.Errorf("URL = %q, want %q", r.URL, tt.wantURL)
			}
		})
	}
}

func TestExtractHTTPRouteRouteFallbackHeuristic(t *testing.T) {
	tests := []struct {
		name       string
		parentRefs []gwv1.ParentReference
		lookup     gatewayLookupFunc
		wantTLS    bool
		wantURL    string
	}{
		{
			name:       "nil lookup, no parentRef signal defaults to http",
			parentRefs: []gwv1.ParentReference{{Name: "gateway"}},
			wantTLS:    false,
			wantURL:    "http://app.example.com",
		},
		{
			name: "nil lookup, sectionName https",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("https")},
			},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "nil lookup, sectionName tls",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("tls-listener")},
			},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "gateway missing from cache, sectionName websecure (traefik convention)",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("websecure")},
			},
			lookup:  lookupFrom(), // never finds anything
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "gateway missing from cache, sectionName insecure stays http",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("web-insecure")},
			},
			lookup:  lookupFrom(),
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
		{
			name: "gateway missing from cache, sectionName web stays http",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("web")},
			},
			lookup:  lookupFrom(),
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
		{
			name: "gateway missing from cache, port 443 without sectionName",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", Port: portNumber(443)},
			},
			lookup:  lookupFrom(),
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "gateway missing from cache, port 80 stays http",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", Port: portNumber(80)},
			},
			lookup:  lookupFrom(),
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
		{
			name: "one of multiple parentRefs is https",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("web")},
				{Name: "gateway", SectionName: sectionName("websecure")},
			},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hr := newHTTPRoute(tt.parentRefs)

			r := extractHTTPRouteRoute(hr, tt.lookup)

			if r.TLS != tt.wantTLS {
				t.Errorf("TLS = %v, want %v", r.TLS, tt.wantTLS)
			}
			if r.URL != tt.wantURL {
				t.Errorf("URL = %q, want %q", r.URL, tt.wantURL)
			}
		})
	}
}

func TestExtractHTTPRouteRouteURLAnnotationOverride(t *testing.T) {
	hr := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "default",
			Annotations: map[string]string{
				"routeboard.io/url": "https://override.example.com",
			},
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{
				ParentRefs: []gwv1.ParentReference{{Name: "gateway"}},
			},
			Hostnames: []gwv1.Hostname{"app.example.com"},
		},
	}

	r := extractHTTPRouteRoute(hr, nil)

	if r.URL != "https://override.example.com" {
		t.Errorf("URL = %q, want annotation override", r.URL)
	}
}
