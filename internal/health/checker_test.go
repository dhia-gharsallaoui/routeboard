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

	status, elapsed := checker.check(&model.Route{URL: srv.URL})
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

	status, elapsed := checker.check(&model.Route{URL: srv.URL})
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

	status, elapsed := checker.check(&model.Route{URL: srv.URL})
	if status != model.HealthDegraded {
		t.Errorf("got %q, want %q", status, model.HealthDegraded)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for server response")
	}
}

func TestCheckUnhealthyConnectionRefused(t *testing.T) {
	checker := NewChecker(testConfig(), store.New(nil))

	status, elapsed := checker.check(&model.Route{URL: "http://localhost:1"}) // nothing listening
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

	status, elapsed := checker.check(&model.Route{URL: redirect.URL})
	if status != model.HealthHealthy {
		t.Errorf("got %q, want %q", status, model.HealthHealthy)
	}
	if elapsed <= 0 {
		t.Error("expected elapsed > 0 for healthy redirect")
	}
}

func TestCheckDefaultsToHeadOnRoot(t *testing.T) {
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewChecker(testConfig(), store.New(nil))

	if status, _ := checker.check(&model.Route{URL: srv.URL}); status != model.HealthHealthy {
		t.Errorf("got %q, want %q", status, model.HealthHealthy)
	}
	if gotMethod != http.MethodHead {
		t.Errorf("method = %q, want %q", gotMethod, http.MethodHead)
	}
	if gotPath != "/" {
		t.Errorf("path = %q, want %q", gotPath, "/")
	}
}

func TestCheckUsesHealthPathAndMethod(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotQuery = r.Method, r.URL.Path, r.URL.RawQuery
		if r.Method != http.MethodGet || r.URL.Path != "/healthz" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewChecker(testConfig(), store.New(nil))

	route := &model.Route{
		URL:          srv.URL + "/app",
		HealthPath:   "/healthz?verbose=1",
		HealthMethod: http.MethodGet,
	}
	if status, _ := checker.check(route); status != model.HealthHealthy {
		t.Errorf("got %q, want %q", status, model.HealthHealthy)
	}
	if gotMethod != http.MethodGet {
		t.Errorf("method = %q, want %q", gotMethod, http.MethodGet)
	}
	if gotPath != "/healthz" {
		t.Errorf("path = %q, want %q", gotPath, "/healthz")
	}
	if gotQuery != "verbose=1" {
		t.Errorf("query = %q, want %q", gotQuery, "verbose=1")
	}
}

func TestHealthURL(t *testing.T) {
	tests := []struct {
		name  string
		route *model.Route
		want  string
	}{
		{
			name:  "no health path keeps the route url",
			route: &model.Route{URL: "https://example.com/app"},
			want:  "https://example.com/app",
		},
		{
			name:  "health path replaces the url path",
			route: &model.Route{URL: "https://example.com/app", HealthPath: "/healthz"},
			want:  "https://example.com/healthz",
		},
		{
			name:  "health path keeps its query",
			route: &model.Route{URL: "https://example.com/app?a=1", HealthPath: "/healthz?b=2"},
			want:  "https://example.com/healthz?b=2",
		},
		{
			name:  "host in the health path is ignored",
			route: &model.Route{URL: "https://example.com/app", HealthPath: "//evil.example.com/healthz"},
			want:  "https://example.com/healthz",
		},
		{
			name:  "unparseable url is left alone",
			route: &model.Route{URL: "://nope", HealthPath: "/healthz"},
			want:  "://nope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := healthURL(tt.route); got != tt.want {
				t.Errorf("healthURL() = %q, want %q", got, tt.want)
			}
		})
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
