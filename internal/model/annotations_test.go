package model

import "testing"

func TestApplyAnnotations(t *testing.T) {
	r := &Route{
		Title: "default-title",
		Group: "default-group",
	}

	annotations := map[string]string{
		"routeboard.io/title":       "Custom Title",
		"routeboard.io/description": "A great service",
		"routeboard.io/icon":        "🎯",
		"routeboard.io/group":       "Monitoring",
		"routeboard.io/order":       "5",
		"routeboard.io/hidden":      "true",
		"routeboard.io/url":         "https://custom.example.com",
		"routeboard.io/health":      "false",
		"unrelated/annotation":      "ignored",
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
