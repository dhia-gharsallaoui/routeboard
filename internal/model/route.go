package model

import "time"

type RouteSource string

const (
	SourceIngress   RouteSource = "Ingress"
	SourceHTTPRoute RouteSource = "HTTPRoute"
)

type Route struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Source      RouteSource       `json:"source"`
	URL         string            `json:"url"`
	Hosts       []string          `json:"hosts"`
	Paths       []string          `json:"paths"`
	TLS         bool              `json:"tls"`
	ServiceName string            `json:"serviceName,omitempty"`
	ServicePort string            `json:"servicePort,omitempty"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Group       string            `json:"group"`
	Order       int               `json:"order"`
	Hidden      bool              `json:"hidden"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"-"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`

	Health              HealthStatus   `json:"health"`
	HealthCheckedAt     time.Time      `json:"healthCheckedAt,omitempty"`
	HealthHistory       []HealthStatus `json:"healthHistory,omitempty"`
	ResponseTimeMs      int64          `json:"responseTimeMs,omitempty"`
	ResponseTimeHistory []int64        `json:"responseTimeHistory,omitempty"`
}

const HealthHistoryMax = 60
