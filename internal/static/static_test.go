package static

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dhia/routeboard/internal/model"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "routes.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadValidFile(t *testing.T) {
	path := writeFile(t, `
routes:
  - name: grafana
    url: http://192.168.10.30:3000
    title: Grafana
    description: Metrics & dashboards
    group: monitoring
    icon: si:grafana
    order: 1
  - name: proxmox
    url: https://192.168.1.10:8006
    health: false
`)
	routes, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(routes) != 2 {
		t.Fatalf("len(routes) = %d, want 2", len(routes))
	}

	g := routes[0]
	if g.ID != "static/grafana" {
		t.Errorf("ID = %q, want %q", g.ID, "static/grafana")
	}
	if g.Source != model.SourceStatic {
		t.Errorf("Source = %q, want %q", g.Source, model.SourceStatic)
	}
	if g.Namespace != "static" {
		t.Errorf("Namespace = %q, want %q", g.Namespace, "static")
	}
	if g.Title != "Grafana" || g.Group != "monitoring" || g.Icon != "si:grafana" || g.Order != 1 {
		t.Errorf("unexpected fields: %+v", g)
	}
	if g.Description != "Metrics & dashboards" {
		t.Errorf("Description = %q", g.Description)
	}
	if g.HealthDisabled {
		t.Error("HealthDisabled = true, want false (default)")
	}
	if g.Health != model.HealthUnknown {
		t.Errorf("Health = %q, want %q", g.Health, model.HealthUnknown)
	}
	if g.CreatedAt.IsZero() || g.UpdatedAt.IsZero() {
		t.Error("CreatedAt/UpdatedAt not set")
	}

	p := routes[1]
	if p.Title != "proxmox" {
		t.Errorf("Title = %q, want name fallback %q", p.Title, "proxmox")
	}
	if p.Group != "static" {
		t.Errorf("Group = %q, want default %q", p.Group, "static")
	}
	if !p.HealthDisabled {
		t.Error("HealthDisabled = false, want true (health: false)")
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{"missing name", "routes:\n  - url: http://x\n", "name is required"},
		{"missing url", "routes:\n  - name: a\n", "url is required"},
		{"duplicate name", "routes:\n  - name: a\n    url: http://x\n  - name: a\n    url: http://y\n", "duplicate name"},
		{"invalid url", "routes:\n  - name: a\n    url: '://nope'\n", "invalid url"},
		{"url without scheme", "routes:\n  - name: a\n    url: 192.168.1.10:8006\n", "invalid url"},
		{"empty file", "routes: []\n", "no routes"},
		{"malformed yaml", "routes: [oops\n", "parse static routes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(writeFile(t, tt.content))
			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q does not contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yml"))
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}
