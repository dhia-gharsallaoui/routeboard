package k8s

import (
	"fmt"
	"strings"
	"time"

	networkingv1 "k8s.io/api/networking/v1"

	"github.com/dhia/routeboard/internal/model"
)

func extractIngressRoute(ingress *networkingv1.Ingress) *model.Route {
	r := &model.Route{
		ID:          fmt.Sprintf("Ingress:%s/%s", ingress.Namespace, ingress.Name),
		Name:        ingress.Name,
		Namespace:   ingress.Namespace,
		Source:      model.SourceIngress,
		Group:       ingress.Namespace,
		Labels:      ingress.Labels,
		Annotations: ingress.Annotations,
		CreatedAt:   ingress.CreationTimestamp.Time,
		UpdatedAt:   time.Now(),
	}

	tlsHosts := make(map[string]bool)
	for _, tls := range ingress.Spec.TLS {
		for _, host := range tls.Hosts {
			tlsHosts[host] = true
		}
	}
	if len(ingress.Spec.TLS) > 0 {
		r.TLS = true
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.Host != "" {
			r.Hosts = append(r.Hosts, rule.Host)
		}
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				if path.Path != "" {
					r.Paths = append(r.Paths, path.Path)
				}
				if r.ServiceName == "" && path.Backend.Service != nil {
					r.ServiceName = path.Backend.Service.Name
					if path.Backend.Service.Port.Name != "" {
						r.ServicePort = path.Backend.Service.Port.Name
					} else {
						r.ServicePort = fmt.Sprintf("%d", path.Backend.Service.Port.Number)
					}
				}
			}
		}
	}

	r.URL = computeURL(r)
	r.Title = titleize(r.Name)
	r.Icon = model.DetectIcon(r.ServiceName, r.Name)
	model.ApplyAnnotations(r, ingress.Annotations)

	return r
}

func computeURL(r *model.Route) string {
	scheme := "http"
	if r.TLS {
		scheme = "https"
	}
	host := ""
	if len(r.Hosts) > 0 {
		host = r.Hosts[0]
	}
	if host == "" {
		return ""
	}
	path := ""
	if len(r.Paths) > 0 && r.Paths[0] != "/" {
		path = r.Paths[0]
	}
	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func titleize(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_' || r == '.'
	})
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
