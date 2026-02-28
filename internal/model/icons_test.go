package model

import "testing"

func TestDetectIconKnownServices(t *testing.T) {
	tests := []struct {
		service string
		want    string
	}{
		{"grafana", "📊"},
		{"my-grafana-instance", "📊"},
		{"prometheus-server", "🔥"},
		{"argocd-server", "🚀"},
		{"argo-cd-server", "🚀"},
		{"vault", "🔐"},
		{"keycloak", "🔑"},
		{"jaeger-query", "🔍"},
		{"harbor-core", "⚓"},
		{"redis-master", "🔴"},
		{"rabbitmq", "🐰"},
	}

	for _, tt := range tests {
		got := DetectIcon(tt.service, "")
		if got != tt.want {
			t.Errorf("DetectIcon(%q) = %q, want %q", tt.service, got, tt.want)
		}
	}
}

func TestDetectIconFallbackToResourceName(t *testing.T) {
	got := DetectIcon("", "grafana-ingress")
	if got != "📊" {
		t.Errorf("DetectIcon(\"\", \"grafana-ingress\") = %q, want 📊", got)
	}
}

func TestDetectIconDefault(t *testing.T) {
	got := DetectIcon("my-custom-app", "my-custom-app")
	if got != "🌐" {
		t.Errorf("DetectIcon(\"my-custom-app\") = %q, want 🌐", got)
	}
}
