package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/model"
	"github.com/dhia/routeboard/internal/store"
)

func testConfig() *config.Config {
	return &config.Config{
		HealthEnabled:  true,
		HealthInterval: time.Minute,
		HealthTimeout:  2 * time.Second,
	}
}

func TestCheckHealthy(t *testing.T) {
	requestedMethod := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedMethod <- r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := testConfig()
	s := store.New(nil)
	checker := NewChecker(cfg, s)

	status, elapsed := checker.check(srv.URL, http.MethodHead)
	if status != model.HealthHealthy {
		t.Errorf("got %q, want %q", status, model.HealthHealthy)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for healthy check")
	}
	if got := <-requestedMethod; got != http.MethodHead {
		t.Errorf("method = %q, want %q", got, http.MethodHead)
	}
}

func TestCheckUnhealthy500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	checker := NewChecker(testConfig(), store.New(nil))

	status, elapsed := checker.check(srv.URL, http.MethodHead)
	if status != model.HealthUnhealthy {
		t.Errorf("got %q, want %q", status, model.HealthUnhealthy)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for server response")
	}
}

func TestCheckDegraded4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	checker := NewChecker(testConfig(), store.New(nil))

	status, elapsed := checker.check(srv.URL, http.MethodHead)
	if status != model.HealthDegraded {
		t.Errorf("got %q, want %q", status, model.HealthDegraded)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for server response")
	}
}

func TestCheckUnhealthyConnectionRefused(t *testing.T) {
	checker := NewChecker(testConfig(), store.New(nil))

	status, elapsed := checker.check("http://localhost:1", http.MethodHead) // nothing listening
	if status != model.HealthUnhealthy {
		t.Errorf("got %q, want %q", status, model.HealthUnhealthy)
	}
	if elapsed != 0 {
		t.Errorf("expected elapsed == 0 for connection error, got %v", elapsed)
	}
}

func TestCheckHealthyRedirect(t *testing.T) {
	final := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer final.Close()

	redirect := httptest.NewServer(http.RedirectHandler(final.URL, http.StatusFound))
	defer redirect.Close()

	checker := NewChecker(testConfig(), store.New(nil))

	status, elapsed := checker.check(redirect.URL, http.MethodHead)
	if status != model.HealthHealthy {
		t.Errorf("got %q, want %q", status, model.HealthHealthy)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for healthy redirect")
	}
}

func TestCheckAllUpdatesStore(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	var events []store.ChangeEvent
	s := store.New(func(e store.ChangeEvent) { events = append(events, e) })
	s.Set(&model.Route{ID: "test-1", URL: srv.URL, Title: "Test"})

	checker := NewChecker(testConfig(), s)
	checker.checkAll()

	route, _ := s.Get("test-1")
	if route.Health != model.HealthHealthy {
		t.Errorf("route health = %q, want %q", route.Health, model.HealthHealthy)
	}
	if route.HealthCheckedAt.IsZero() {
		t.Error("HealthCheckedAt should be set")
	}
	if route.ResponseTimeMs < 0 {
		t.Error("expected ResponseTimeMs >= 0 after health check")
	}
	if len(route.ResponseTimeHistory) != 1 {
		t.Errorf("expected 1 response time history entry, got %d", len(route.ResponseTimeHistory))
	}

	hasHealthEvent := false
	for _, e := range events {
		if e.Type == store.ChangeHealth {
			hasHealthEvent = true
		}
	}
	if !hasHealthEvent {
		t.Error("expected a ChangeHealth event")
	}
}

func TestCheckAllUsesHealthEndpoint(t *testing.T) {
	requestedPath := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Errorf("method = %q, want HEAD", r.Method)
		}
		requestedPath <- r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := store.New(nil)
	s.Set(&model.Route{
		ID:        "health-endpoint",
		URL:       srv.URL + "/dashboard",
		HealthURL: srv.URL + "/ping",
		Title:     "Health endpoint override",
	})

	checker := NewChecker(testConfig(), s)
	checker.checkAll()

	if got := <-requestedPath; got != "/ping" {
		t.Errorf("health request path = %q, want %q", got, "/ping")
	}

	route, _ := s.Get("health-endpoint")
	if route.URL != srv.URL+"/dashboard" {
		t.Errorf("dashboard URL = %q, want unchanged URL %q", route.URL, srv.URL+"/dashboard")
	}
	if route.Health != model.HealthHealthy {
		t.Errorf("route health = %q, want %q", route.Health, model.HealthHealthy)
	}
}

func TestCheckAllUsesHealthMethod(t *testing.T) {
	requestedMethod := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedMethod <- r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := store.New(nil)
	s.Set(&model.Route{
		ID:           "get-health",
		URL:          srv.URL,
		HealthMethod: http.MethodGet,
	})

	checker := NewChecker(testConfig(), s)
	checker.checkAll()

	if got := <-requestedMethod; got != http.MethodGet {
		t.Errorf("method = %q, want %q", got, http.MethodGet)
	}
}

func TestCheckAllSkipsHealthDisabled(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := store.New(nil)
	s.Set(&model.Route{ID: "disabled", URL: srv.URL, Title: "Disabled", HealthDisabled: true})

	checker := NewChecker(testConfig(), s)
	checker.checkAll()

	if requests != 0 {
		t.Errorf("expected 0 requests to health-disabled route, got %d", requests)
	}
	route, _ := s.Get("disabled")
	if route.Health != "" {
		t.Errorf("route health = %q, want empty (unchanged)", route.Health)
	}
}

func TestCheckAllSkipsNoURL(t *testing.T) {
	s := store.New(nil)
	s.Set(&model.Route{ID: "no-url", URL: "", Title: "No URL"})

	checker := NewChecker(testConfig(), s)
	checker.checkAll()

	route, _ := s.Get("no-url")
	if route.Health != "" {
		t.Errorf("route health = %q, want empty (unchanged)", route.Health)
	}
}
