package model

import "testing"

func TestApplyAnnotations(t *testing.T) {
	r := &Route{
		Title: "default-title",
		Group: "default-group",
	}

	annotations := map[string]string{
		"routeboard.io/title":         "Custom Title",
		"routeboard.io/description":   "A great service",
		"routeboard.io/icon":          "🎯",
		"routeboard.io/group":         "Monitoring",
		"routeboard.io/order":         "5",
		"routeboard.io/hidden":        "true",
		"routeboard.io/url":           "https://custom.example.com",
		"routeboard.io/health":        "false",
		"routeboard.io/health-path":   "/healthz",
		"routeboard.io/health-method": "get",
		"unrelated/annotation":        "ignored",
	}

	ApplyAnnotations(r, annotations)

	if r.Title != "Custom Title" {
		t.Errorf("Title = %q, want %q", r.Title, "Custom Title")
	}
	if r.Description != "A great service" {
		t.Errorf("Description = %q, want %q", r.Description, "A great service")
	}
	if r.Icon != "🎯" {
		t.Errorf("Icon = %q, want %q", r.Icon, "🎯")
	}
	if r.Group != "Monitoring" {
		t.Errorf("Group = %q, want %q", r.Group, "Monitoring")
	}
	if r.Order != 5 {
		t.Errorf("Order = %d, want %d", r.Order, 5)
	}
	if !r.Hidden {
		t.Error("Hidden = false, want true")
	}
	if r.URL != "https://custom.example.com" {
		t.Errorf("URL = %q, want %q", r.URL, "https://custom.example.com")
	}
	if !r.HealthDisabled {
		t.Error("HealthDisabled = false, want true")
	}
	if r.HealthPath != "/healthz" {
		t.Errorf("HealthPath = %q, want %q", r.HealthPath, "/healthz")
	}
	if r.HealthMethod != "GET" {
		t.Errorf("HealthMethod = %q, want %q", r.HealthMethod, "GET")
	}
}

func TestApplyAnnotationsHealthPath(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"/healthz", "/healthz"},
		{"healthz", "/healthz"},
		{"  /healthz  ", "/healthz"},
		{"/healthz?verbose=1", "/healthz?verbose=1"},
		{"", "/default"},
		{"   ", "/default"},
	}

	for _, tt := range tests {
		r := &Route{HealthPath: "/default"}
		ApplyAnnotations(r, map[string]string{"routeboard.io/health-path": tt.value})
		if r.HealthPath != tt.want {
			t.Errorf("health-path=%q: HealthPath = %q, want %q", tt.value, r.HealthPath, tt.want)
		}
	}
}

func TestApplyAnnotationsHealthMethod(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"GET", "GET"},
		{"get", "GET"},
		{" head ", "HEAD"},
		{"options", "OPTIONS"},
		// Unsupported or invalid methods keep whatever was already set.
		{"POST", "HEAD"},
		{"DELETE", "HEAD"},
		{"not a method", "HEAD"},
		{"", "HEAD"},
	}

	for _, tt := range tests {
		r := &Route{HealthMethod: "HEAD"}
		ApplyAnnotations(r, map[string]string{"routeboard.io/health-method": tt.value})
		if r.HealthMethod != tt.want {
			t.Errorf("health-method=%q: HealthMethod = %q, want %q", tt.value, r.HealthMethod, tt.want)
		}
	}
}

func TestApplyAnnotationsHealthDefaultsUnset(t *testing.T) {
	r := &Route{}
	ApplyAnnotations(r, map[string]string{"routeboard.io/title": "Grafana"})
	if r.HealthPath != "" {
		t.Errorf("HealthPath = %q, want empty", r.HealthPath)
	}
	if r.HealthMethod != "" {
		t.Errorf("HealthMethod = %q, want empty", r.HealthMethod)
	}
}

func TestApplyAnnotationsHealthNotFalse(t *testing.T) {
	for _, v := range []string{"true", "", "yes", "0"} {
		r := &Route{}
		ApplyAnnotations(r, map[string]string{
			"routeboard.io/health": v,
		})
		if r.HealthDisabled {
			t.Errorf("HealthDisabled = true for health=%q, want false", v)
		}
	}
}

func TestApplyAnnotationsInvalidOrder(t *testing.T) {
	r := &Route{Order: 99}
	ApplyAnnotations(r, map[string]string{
		"routeboard.io/order": "not-a-number",
	})
	if r.Order != 99 {
		t.Errorf("Order = %d, want 99 (unchanged)", r.Order)
	}
}

func TestApplyAnnotationsEmpty(t *testing.T) {
	r := &Route{Title: "Original"}
	ApplyAnnotations(r, map[string]string{})
	if r.Title != "Original" {
		t.Errorf("Title = %q, want %q (unchanged)", r.Title, "Original")
	}
}
