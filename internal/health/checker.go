package health

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
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

	// Filter to routes with URLs
	var targets []*model.Route
	for _, r := range routes {
		if r.URL != "" {
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

			status := c.check(r.URL)
			now := time.Now()
			c.store.UpdateHealth(r.ID, status, now)
		}(route)
	}

	wg.Wait()
}

func (c *Checker) check(url string) model.HealthStatus {
	start := time.Now()

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return model.HealthUnhealthy
	}
	req.Header.Set("User-Agent", "RouteBoard-HealthCheck/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return model.HealthUnhealthy
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)
	degradedThreshold := c.cfg.HealthTimeout / 2

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 400:
		if elapsed > degradedThreshold {
			return model.HealthDegraded
		}
		return model.HealthHealthy
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		return model.HealthDegraded
	default:
		return model.HealthUnhealthy
	}
}
