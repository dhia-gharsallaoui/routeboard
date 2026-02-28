package store

import (
	"sort"
	"sync"
	"time"

	"github.com/dhia/routeboard/internal/model"
)

type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeUpdated ChangeType = "updated"
	ChangeDeleted ChangeType = "deleted"
	ChangeHealth  ChangeType = "health"
)

type ChangeEvent struct {
	Type  ChangeType   `json:"type"`
	Route *model.Route `json:"route"`
}

type NotifyFunc func(event ChangeEvent)

type Store struct {
	mu       sync.RWMutex
	routes   map[string]*model.Route
	notifyFn NotifyFunc
}

func New(notifyFn NotifyFunc) *Store {
	return &Store{
		routes:   make(map[string]*model.Route),
		notifyFn: notifyFn,
	}
}

func (s *Store) Set(route *model.Route) {
	s.mu.Lock()
	existing, exists := s.routes[route.ID]
	s.routes[route.ID] = route
	s.mu.Unlock()

	if s.notifyFn != nil {
		if exists && routeChanged(existing, route) {
			s.notifyFn(ChangeEvent{Type: ChangeUpdated, Route: route})
		} else if !exists {
			s.notifyFn(ChangeEvent{Type: ChangeAdded, Route: route})
		}
	}
}

func (s *Store) Delete(id string) {
	s.mu.Lock()
	route, exists := s.routes[id]
	delete(s.routes, id)
	s.mu.Unlock()

	if exists && s.notifyFn != nil {
		s.notifyFn(ChangeEvent{Type: ChangeDeleted, Route: route})
	}
}

func (s *Store) Get(id string) (*model.Route, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.routes[id]
	return r, ok
}

func (s *Store) List() []*model.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()

	routes := make([]*model.Route, 0, len(s.routes))
	for _, r := range s.routes {
		if !r.Hidden {
			routes = append(routes, r)
		}
	}

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Group != routes[j].Group {
			return routes[i].Group < routes[j].Group
		}
		if routes[i].Order != routes[j].Order {
			return routes[i].Order < routes[j].Order
		}
		return routes[i].Title < routes[j].Title
	})

	return routes
}

func (s *Store) UpdateHealth(id string, status model.HealthStatus, checkedAt time.Time) {
	s.mu.Lock()
	route, exists := s.routes[id]
	if exists {
		route.Health = status
		route.HealthCheckedAt = checkedAt
	}
	s.mu.Unlock()

	if exists && s.notifyFn != nil {
		s.notifyFn(ChangeEvent{Type: ChangeHealth, Route: route})
	}
}

func (s *Store) ListAll() []*model.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	routes := make([]*model.Route, 0, len(s.routes))
	for _, r := range s.routes {
		routes = append(routes, r)
	}
	return routes
}

func (s *Store) GroupedRoutes() map[string][]*model.Route {
	routes := s.List()
	grouped := make(map[string][]*model.Route)
	for _, r := range routes {
		grouped[r.Group] = append(grouped[r.Group], r)
	}
	return grouped
}

func (s *Store) Namespaces() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seen := make(map[string]struct{})
	for _, r := range s.routes {
		seen[r.Namespace] = struct{}{}
	}

	ns := make([]string, 0, len(seen))
	for n := range seen {
		ns = append(ns, n)
	}
	sort.Strings(ns)
	return ns
}

func routeChanged(old, new *model.Route) bool {
	return old.URL != new.URL ||
		old.Title != new.Title ||
		old.Description != new.Description ||
		old.Icon != new.Icon ||
		old.Group != new.Group ||
		old.Order != new.Order ||
		old.Hidden != new.Hidden ||
		old.TLS != new.TLS
}
