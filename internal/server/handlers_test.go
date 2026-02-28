package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/model"
	"github.com/dhia/routeboard/internal/store"
)

func testStore() *store.Store {
	s := store.New(nil)
	s.Set(&model.Route{
		ID: "r1", Name: "grafana", Namespace: "monitoring",
		Title: "Grafana", URL: "https://grafana.example.com",
		Source: model.SourceIngress, Group: "monitoring",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	s.Set(&model.Route{
		ID: "r2", Name: "vault", Namespace: "security",
		Title: "Vault", URL: "https://vault.example.com",
		Source: model.SourceHTTPRoute, Group: "security",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	return s
}

func TestAPIRoutes(t *testing.T) {
	h := NewHandlers(&config.Config{Title: "Test"}, testStore())

	req := httptest.NewRequest("GET", "/api/routes", nil)
	w := httptest.NewRecorder()
	h.APIRoutes(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var routes []model.Route
	json.NewDecoder(w.Body).Decode(&routes)
	if len(routes) != 2 {
		t.Fatalf("got %d routes, want 2", len(routes))
	}
}

func TestAPIRoutesFilterNamespace(t *testing.T) {
	h := NewHandlers(&config.Config{}, testStore())

	req := httptest.NewRequest("GET", "/api/routes?namespace=monitoring", nil)
	w := httptest.NewRecorder()
	h.APIRoutes(w, req)

	var routes []model.Route
	json.NewDecoder(w.Body).Decode(&routes)
	if len(routes) != 1 {
		t.Fatalf("got %d routes, want 1", len(routes))
	}
	if routes[0].Name != "grafana" {
		t.Errorf("got %q, want grafana", routes[0].Name)
	}
}

func TestAPIRoutesFilterSearch(t *testing.T) {
	h := NewHandlers(&config.Config{}, testStore())

	req := httptest.NewRequest("GET", "/api/routes?q=vault", nil)
	w := httptest.NewRecorder()
	h.APIRoutes(w, req)

	var routes []model.Route
	json.NewDecoder(w.Body).Decode(&routes)
	if len(routes) != 1 {
		t.Fatalf("got %d routes, want 1", len(routes))
	}
	if routes[0].Name != "vault" {
		t.Errorf("got %q, want vault", routes[0].Name)
	}
}

func TestAPIConfig(t *testing.T) {
	h := NewHandlers(&config.Config{Title: "My Board"}, testStore())

	req := httptest.NewRequest("GET", "/api/config", nil)
	w := httptest.NewRecorder()
	h.APIConfig(w, req)

	var resp configResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Title != "My Board" {
		t.Errorf("Title = %q, want %q", resp.Title, "My Board")
	}
	if len(resp.Namespaces) != 2 {
		t.Errorf("got %d namespaces, want 2", len(resp.Namespaces))
	}
}

func TestHealth(t *testing.T) {
	h := NewHandlers(&config.Config{}, store.New(nil))

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("body = %q, want %q", w.Body.String(), "ok")
	}
}
