package model

import (
	"net/url"
	"strconv"
)

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
	if v, ok := annotations[AnnotationPrefix+"health"]; ok {
		r.HealthDisabled = v == "false"
	}
	if v, ok := annotations[AnnotationPrefix+"health-endpoint"]; ok {
		r.HealthURL = resolveHealthURL(r.URL, v)
	}
}

// resolveHealthURL combines a health endpoint with the route origin while
// leaving the dashboard URL unchanged. An absolute endpoint URL is used as-is.
func resolveHealthURL(baseURL, endpoint string) string {
	endpointURL, err := url.Parse(endpoint)
	if err != nil || endpointURL.String() == "" {
		return ""
	}
	if endpointURL.IsAbs() {
		if endpointURL.Host == "" {
			return ""
		}
		return endpointURL.String()
	}

	base, err := url.Parse(baseURL)
	if err != nil || base.Scheme == "" || base.Host == "" {
		return ""
	}
	if endpointURL.Host != "" {
		return ""
	}

	base.Path = endpointURL.Path
	base.RawPath = endpointURL.RawPath
	base.RawQuery = endpointURL.RawQuery
	base.Fragment = ""
	return base.String()
}
