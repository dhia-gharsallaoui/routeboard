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
