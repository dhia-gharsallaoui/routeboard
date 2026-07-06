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

func TestExtractHTTPRouteRouteScheme(t *testing.T) {
	tests := []struct {
		name       string
		parentRefs []gwv1.ParentReference
		wantTLS    bool
		wantURL    string
	}{
		{
			name:       "no parentRef signal defaults to http",
			parentRefs: []gwv1.ParentReference{{Name: "gateway"}},
			wantTLS:    false,
			wantURL:    "http://app.example.com",
		},
		{
			name: "sectionName https",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("https")},
			},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "sectionName tls",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("tls-listener")},
			},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "sectionName websecure (traefik convention)",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("websecure")},
			},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "sectionName insecure stays http",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("web-insecure")},
			},
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
		{
			name: "sectionName web stays http",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", SectionName: sectionName("web")},
			},
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
		{
			name: "port 443 without sectionName",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", Port: portNumber(443)},
			},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
		{
			name: "port 80 stays http",
			parentRefs: []gwv1.ParentReference{
				{Name: "gateway", Port: portNumber(80)},
			},
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
			hr := &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "default"},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: tt.parentRefs},
					Hostnames:       []gwv1.Hostname{"app.example.com"},
				},
			}

			r := extractHTTPRouteRoute(hr)

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

	r := extractHTTPRouteRoute(hr)

	if r.URL != "https://override.example.com" {
		t.Errorf("URL = %q, want annotation override", r.URL)
	}
}
