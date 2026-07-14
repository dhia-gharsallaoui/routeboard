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

	KubeEnabled      bool
	StaticRoutesPath string

	Title     string
	LogLevel  string
	LogFormat string // text, json

	HealthEnabled  bool
	HealthInterval time.Duration
	HealthTimeout  time.Duration

	WebhookURL    string
	WebhookFormat string // json, slack, discord
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
		LabelSelector:    envStr("ROUTEBOARD_LABEL_SELECTOR", ""),
		WatchIngress:     envBool("ROUTEBOARD_WATCH_INGRESS", true),
		WatchHTTPRoute:   envBool("ROUTEBOARD_WATCH_HTTPROUTE", true),
		KubeEnabled:      envBool("ROUTEBOARD_KUBE_ENABLED", true),
		StaticRoutesPath: envStr("ROUTEBOARD_STATIC_ROUTES", ""),
		Title:            envStr("ROUTEBOARD_TITLE", "RouteBoard"),
		LogLevel:         envStr("ROUTEBOARD_LOG_LEVEL", "info"),
		LogFormat:        envStr("ROUTEBOARD_LOG_FORMAT", "text"),
		HealthEnabled:    envBool("ROUTEBOARD_HEALTH_ENABLED", true),
		HealthInterval:   envDuration("ROUTEBOARD_HEALTH_INTERVAL", 30*time.Second),
		HealthTimeout:    envDuration("ROUTEBOARD_HEALTH_TIMEOUT", 5*time.Second),
		WebhookURL:       envStr("ROUTEBOARD_WEBHOOK_URL", ""),
		WebhookFormat:    envStr("ROUTEBOARD_WEBHOOK_FORMAT", "json"),
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
