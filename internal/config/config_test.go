package config

import "testing"

func TestLoadLogFormat(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want string
	}{
		{"default is text", "", "text"},
		{"json accepted", "json", "json"},
		{"text accepted", "text", "text"},
		{"invalid value passed through for fallback at setup", "yaml", "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				t.Setenv("ROUTEBOARD_LOG_FORMAT", tt.env)
			}
			cfg := Load()
			if cfg.LogFormat != tt.want {
				t.Errorf("LogFormat = %q; want %q", cfg.LogFormat, tt.want)
			}
		})
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.StaticRoutesPath != "" {
		t.Errorf("StaticRoutesPath = %q, want empty default", cfg.StaticRoutesPath)
	}
	if !cfg.KubeEnabled {
		t.Error("KubeEnabled = false, want true by default")
	}
}

func TestLoadStaticAndKubeEnv(t *testing.T) {
	t.Setenv("ROUTEBOARD_STATIC_ROUTES", "/etc/routeboard/routes.yml")
	t.Setenv("ROUTEBOARD_KUBE_ENABLED", "false")

	cfg := Load()
	if cfg.StaticRoutesPath != "/etc/routeboard/routes.yml" {
		t.Errorf("StaticRoutesPath = %q, want %q", cfg.StaticRoutesPath, "/etc/routeboard/routes.yml")
	}
	if cfg.KubeEnabled {
		t.Error("KubeEnabled = true, want false")
	}
}
