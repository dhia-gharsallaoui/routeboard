package k8s

import (
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractIngressRouteScheme(t *testing.T) {
	tests := []struct {
		name    string
		tls     []networkingv1.IngressTLS
		wantTLS bool
		wantURL string
	}{
		{
			name:    "no tls defaults to http",
			tls:     nil,
			wantTLS: false,
			wantURL: "http://app.example.com",
		},
		{
			name:    "tls covering host",
			tls:     []networkingv1.IngressTLS{{Hosts: []string{"app.example.com"}}},
			wantTLS: true,
			wantURL: "https://app.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingress := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "default"},
				Spec: networkingv1.IngressSpec{
					TLS: tt.tls,
					Rules: []networkingv1.IngressRule{
						{Host: "app.example.com"},
					},
				},
			}

			r := extractIngressRoute(ingress)

			if r.TLS != tt.wantTLS {
				t.Errorf("TLS = %v, want %v", r.TLS, tt.wantTLS)
			}
			if r.URL != tt.wantURL {
				t.Errorf("URL = %q, want %q", r.URL, tt.wantURL)
			}
		})
	}
}
