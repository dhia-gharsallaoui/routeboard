package health

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
	tlsCfg, err := buildHealthTLSConfig(cfg)
	if err != nil {
		slog.Warn("health checker TLS config could not be loaded; continuing without mTLS",
			"error", err)
		tlsCfg = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	return &Checker{
		cfg:   cfg,
		store: s,
		client: &http.Client{
			Timeout: cfg.HealthTimeout,
			Transport: &http.Transport{
				TLSClientConfig: tlsCfg,
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

			status, elapsed := c.check(r.URL)
			now := time.Now()
			c.store.UpdateHealth(r.ID, status, now, elapsed.Milliseconds())
		}(route)
	}

	wg.Wait()
}

func (c *Checker) check(url string) (model.HealthStatus, time.Duration) {
	start := time.Now()

	req, err := http.NewRequest(http.MethodHead, url, nil)
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

func buildHealthTLSConfig(cfg *config.Config) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if cfg.HealthCAFile != "" {
		roots, err := x509.SystemCertPool()
		if err != nil || roots == nil {
			roots = x509.NewCertPool()
		}
		caPEM, err := os.ReadFile(cfg.HealthCAFile)
		if err != nil {
			return nil, fmt.Errorf("read health-check CA: %w", err)
		}
		if ok := roots.AppendCertsFromPEM(caPEM); !ok {
			return nil, fmt.Errorf("health-check CA file contains no valid certificates")
		}
		tlsCfg.RootCAs = roots
	}

	if cfg.HealthClientCert != "" || cfg.HealthClientKey != "" {
		if cfg.HealthClientCert == "" || cfg.HealthClientKey == "" {
			return nil, fmt.Errorf("both client cert and key required")
		}
		cert, err := tls.LoadX509KeyPair(cfg.HealthClientCert, cfg.HealthClientKey)
		if err != nil {
			return nil, err
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	return tlsCfg, nil
}
