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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := testConfig()
	s := store.New(nil)
	checker := NewChecker(cfg, s)

	status, elapsed := checker.check(srv.URL)
	if status != model.HealthHealthy {
		t.Errorf("got %q, want %q", status, model.HealthHealthy)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for healthy check")
	}
}

func TestCheckUnhealthy500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	checker := NewChecker(testConfig(), store.New(nil))

	status, elapsed := checker.check(srv.URL)
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

	status, elapsed := checker.check(srv.URL)
	if status != model.HealthDegraded {
		t.Errorf("got %q, want %q", status, model.HealthDegraded)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for server response")
	}
}

func TestCheckUnhealthyConnectionRefused(t *testing.T) {
	checker := NewChecker(testConfig(), store.New(nil))

	status, elapsed := checker.check("http://localhost:1") // nothing listening
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

	status, elapsed := checker.check(redirect.URL)
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

	// Should have received a health change event
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
