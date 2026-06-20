package model

import "strconv"

const AnnotationPrefix = "routeboard.io/"

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
	// routeboard.io/health: "false" keeps the route listed but disables health checks.
	if v, ok := annotations[AnnotationPrefix+"health"]; ok {
		if enabled, err := strconv.ParseBool(v); err == nil {
			r.MonitorDisabled = !enabled
		}
	}
}
