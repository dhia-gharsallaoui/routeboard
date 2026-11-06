package model

import "testing"

func TestApplyAnnotations(t *testing.T) {
	r := &Route{
		Title: "default-title",
		Group: "default-group",
	}

	annotations := map[string]string{
		"routeboard.io/title":          "Custom Title",
		"routeboard.io/description":    "A great service",
		"routeboard.io/icon":           "🎯",
		"routeboard.io/group":          "Monitoring",
		"routeboard.io/order":          "5",
		"routeboard.io/hidden":         "true",
		"routeboard.io/url":            "https://custom.example.com",
		"routeboard.io/health":         "false",
		"routeboard.io/health-endpoint": "/ping",
		"routeboard.io/health-method":  "GET",
		"unrelated/annotation":         "ignored",
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
	if r.HealthURL != "https://custom.example.com/ping" {
		t.Errorf("HealthURL = %q, want %q", r.HealthURL, "https://custom.example.com/ping")
	}
	if r.HealthMethod != "GET" {
		t.Errorf("HealthMethod = %q, want %q", r.HealthMethod, "GET")
	}
	if !r.HealthDisabled {
		t.Error("HealthDisabled = false, want true")
	}
}

func TestApplyAnnotationsHealthEndpointKeepsDetectedURL(t *testing.T) {
	r := &Route{URL: "https://app.example.com/dashboard"}
	ApplyAnnotations(r, map[string]string{
		"routeboard.io/health-endpoint": "/ping",
	})

	if r.URL != "https://app.example.com/dashboard" {
		t.Errorf("URL = %q, want the detected dashboard URL unchanged", r.URL)
	}
	if r.HealthURL != "https://app.example.com/ping" {
		t.Errorf("HealthURL = %q, want %q", r.HealthURL, "https://app.example.com/ping")
	}
}

func TestApplyAnnotationsHealthMethod(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  string
	}{
		{value: "GET", want: "GET"},
		{value: " head ", want: "HEAD"},
	} {
		r := &Route{}
		ApplyAnnotations(r, map[string]string{
			"routeboard.io/health-method": tc.value,
		})
		if r.HealthMethod != tc.want {
			t.Errorf("HealthMethod for %q = %q, want %q", tc.value, r.HealthMethod, tc.want)
		}
	}
}

func TestApplyAnnotationsInvalidHealthMethodIgnored(t *testing.T) {
	r := &Route{HealthMethod: "GET"}
	ApplyAnnotations(r, map[string]string{
		"routeboard.io/health-method": "POST",
	})
	if r.HealthMethod != "GET" {
		t.Errorf("HealthMethod = %q, want existing value %q", r.HealthMethod, "GET")
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
