// Package static loads dashboard routes from a YAML file, for endpoints
// that don't live behind a Kubernetes Ingress/HTTPRoute (bare VMs, LXCs,
// appliances).
package static

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/dhia/routeboard/internal/model"
)

type entry struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Group       string `yaml:"group"`
	Icon        string `yaml:"icon"`
	Order       int    `yaml:"order"`
	Health      *bool  `yaml:"health"`
}

type file struct {
	Routes []entry `yaml:"routes"`
}

// Load reads and validates a static routes file. Any invalid entry is a
// hard error: the caller is expected to treat it as fatal at startup.
func Load(path string) ([]*model.Route, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read static routes: %w", err)
	}
	var f file
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse static routes: %w", err)
	}
	if len(f.Routes) == 0 {
		return nil, fmt.Errorf("static routes file %s contains no routes", path)
	}

	seen := make(map[string]struct{}, len(f.Routes))
	routes := make([]*model.Route, 0, len(f.Routes))
	now := time.Now()

	for i, e := range f.Routes {
		if e.Name == "" {
			return nil, fmt.Errorf("static route #%d: name is required", i+1)
		}
		if _, dup := seen[e.Name]; dup {
			return nil, fmt.Errorf("static route %q: duplicate name", e.Name)
		}
		seen[e.Name] = struct{}{}
		if e.URL == "" {
			return nil, fmt.Errorf("static route %q: url is required", e.Name)
		}
		if u, err := url.Parse(e.URL); err != nil || u.Scheme == "" || u.Host == "" {
			return nil, fmt.Errorf("static route %q: invalid url %q", e.Name, e.URL)
		}

		title := e.Title
		if title == "" {
			title = e.Name
		}
		group := e.Group
		if group == "" {
			group = "static"
		}

		routes = append(routes, &model.Route{
			ID:             "static/" + e.Name,
			Name:           e.Name,
			Namespace:      "static",
			Source:         model.SourceStatic,
			URL:            e.URL,
			Title:          title,
			Description:    e.Description,
			Icon:           e.Icon,
			Group:          group,
			Order:          e.Order,
			Health:         model.HealthUnknown,
			HealthDisabled: e.Health != nil && !*e.Health,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}
	return routes, nil
}
