package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port    int
	Address string

	Kubeconfig     string
	ResyncInterval time.Duration

	NamespaceAllowlist []string
	NamespaceDenylist  []string
	LabelSelector      string

	WatchIngress   bool
	WatchHTTPRoute bool

	Title    string
	LogLevel string

	HealthEnabled  bool
	HealthInterval time.Duration
	HealthTimeout  time.Duration
}

func Load() *Config {
	return &Config{
		Port:               envInt("ROUTEBOARD_PORT", 8080),
		Address:            envStr("ROUTEBOARD_ADDRESS", "0.0.0.0"),
		Kubeconfig:         envStr("KUBECONFIG", ""),
		ResyncInterval:     envDuration("ROUTEBOARD_RESYNC_INTERVAL", 30*time.Minute),
		NamespaceAllowlist: envSlice("ROUTEBOARD_NAMESPACE_ALLOWLIST", nil),
		NamespaceDenylist: envSlice("ROUTEBOARD_NAMESPACE_DENYLIST", []string{
			"kube-system", "kube-public", "kube-node-lease",
		}),
		LabelSelector:  envStr("ROUTEBOARD_LABEL_SELECTOR", ""),
		WatchIngress:   envBool("ROUTEBOARD_WATCH_INGRESS", true),
		WatchHTTPRoute: envBool("ROUTEBOARD_WATCH_HTTPROUTE", true),
		Title:          envStr("ROUTEBOARD_TITLE", "RouteBoard"),
		LogLevel:       envStr("ROUTEBOARD_LOG_LEVEL", "info"),
		HealthEnabled:  envBool("ROUTEBOARD_HEALTH_ENABLED", true),
		HealthInterval: envDuration("ROUTEBOARD_HEALTH_INTERVAL", 30*time.Second),
		HealthTimeout:  envDuration("ROUTEBOARD_HEALTH_TIMEOUT", 5*time.Second),
	}
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func envSlice(key string, fallback []string) []string {
	if v := os.Getenv(key); v != "" {
		parts := strings.Split(v, ",")
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return fallback
}
