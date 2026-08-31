package health

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/model"
	"github.com/dhia/routeboard/internal/store"
)

type Checker struct {
	cfg    *config.Config
	store  *store.Store
	client *http.Client
}

func NewChecker(cfg *config.Config, s *store.Store) *Checker {
	return &Checker{
		cfg:   cfg,
		store: s,
		client: &http.Client{
			Timeout: cfg.HealthTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (c *Checker) Run(ctx context.Context) {
	slog.Info("starting health checker", "interval", c.cfg.HealthInterval, "timeout", c.cfg.HealthTimeout)

	// Run an initial check after a short delay (let informers populate first)
	select {
	case <-time.After(5 * time.Second):
		c.checkAll()
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(c.cfg.HealthInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("health checker stopped")
			return
		case <-ticker.C:
			c.checkAll()
		}
	}
}

func (c *Checker) checkAll() {
	routes := c.store.ListAll()

	// Filter to routes with URLs and health checks enabled
	var targets []*model.Route
	for _, r := range routes {
		if r.URL != "" && !r.HealthDisabled {
			targets = append(targets, r)
		}
	}

	if len(targets) == 0 {
		return
	}

	slog.Debug("running health checks", "count", len(targets))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // max 10 concurrent checks

	for _, route := range targets {
		wg.Add(1)
		sem <- struct{}{}
		go func(r *model.Route) {
			defer wg.Done()
			defer func() { <-sem }()

			status, elapsed := c.check(r)
			now := time.Now()
			c.store.UpdateHealth(r.ID, status, now, elapsed.Milliseconds())
		}(route)
	}

	wg.Wait()
}

func (c *Checker) check(r *model.Route) (model.HealthStatus, time.Duration) {
	start := time.Now()

	method := r.HealthMethod
	if method == "" {
		method = model.DefaultHealthMethod
	}

	req, err := http.NewRequest(method, healthURL(r), nil)
	if err != nil {
		return model.HealthUnhealthy, 0
	}
	req.Header.Set("User-Agent", "RouteBoard-HealthCheck/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return model.HealthUnhealthy, 0
	}
	defer func() { _ = resp.Body.Close() }()

	elapsed := time.Since(start)
	degradedThreshold := c.cfg.HealthTimeout / 2

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 400:
		if elapsed > degradedThreshold {
			return model.HealthDegraded, elapsed
		}
		return model.HealthHealthy, elapsed
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		return model.HealthDegraded, elapsed
	default:
		return model.HealthUnhealthy, elapsed
	}
}

// healthURL applies the route's health path to its URL. The path is always
// root-relative, so it replaces the URL path instead of extending it; an empty
// health path leaves the route URL untouched.
func healthURL(r *model.Route) string {
	if r.HealthPath == "" {
		return r.URL
	}

	base, err := url.Parse(r.URL)
	if err != nil {
		return r.URL
	}
	ref, err := url.Parse(r.HealthPath)
	if err != nil {
		return r.URL
	}

	// Only the path and query are taken from the annotation - the route URL
	// stays the authority on scheme and host.
	base.Path = ref.Path
	base.RawPath = ref.RawPath
	base.RawQuery = ref.RawQuery
	base.Fragment = ""
	base.RawFragment = ""

	return base.String()
}
