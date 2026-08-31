package model

import (
	"net/http"
	"strconv"
	"strings"
)

const AnnotationPrefix = "routeboard.io/"

// DefaultHealthMethod is used when a route does not set routeboard.io/health-method.
const DefaultHealthMethod = http.MethodHead

// healthMethods are the methods accepted for health probes. Health checks must
// stay side-effect free, so mutating methods are rejected.
var healthMethods = map[string]bool{
	http.MethodHead:    true,
	http.MethodGet:     true,
	http.MethodOptions: true,
}

func ApplyAnnotations(r *Route, annotations map[string]string) {
	if v, ok := annotations[AnnotationPrefix+"title"]; ok {
		r.Title = v
	}
	if v, ok := annotations[AnnotationPrefix+"description"]; ok {
		r.Description = v
	}
	if v, ok := annotations[AnnotationPrefix+"icon"]; ok {
		r.Icon = v
	}
	if v, ok := annotations[AnnotationPrefix+"group"]; ok {
		r.Group = v
	}
	if v, ok := annotations[AnnotationPrefix+"order"]; ok {
		if order, err := strconv.Atoi(v); err == nil {
			r.Order = order
		}
	}
	if v, ok := annotations[AnnotationPrefix+"hidden"]; ok {
		r.Hidden = v == "true"
	}
	if v, ok := annotations[AnnotationPrefix+"url"]; ok {
		r.URL = v
	}
	if v, ok := annotations[AnnotationPrefix+"health"]; ok {
		r.HealthDisabled = v == "false"
	}
	if v, ok := annotations[AnnotationPrefix+"health-path"]; ok {
		if path := NormalizeHealthPath(v); path != "" {
			r.HealthPath = path
		}
	}
	if v, ok := annotations[AnnotationPrefix+"health-method"]; ok {
		if method := NormalizeHealthMethod(v); method != "" {
			r.HealthMethod = method
		}
	}
}

// NormalizeHealthPath trims the value and makes it root-relative. It returns an
// empty string for values that carry no path, so the caller keeps its default.
func NormalizeHealthPath(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "/") {
		v = "/" + v
	}
	return v
}

// NormalizeHealthMethod upper-cases the value and returns an empty string when
// it is not a method allowed for health probes, so the caller keeps its default.
func NormalizeHealthMethod(v string) string {
	v = strings.ToUpper(strings.TrimSpace(v))
	if !healthMethods[v] {
		return ""
	}
	return v
}
