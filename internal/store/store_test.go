package store

import (
	"testing"
	"time"

	"github.com/dhia/routeboard/internal/model"
)

func newRoute(id, title, group string, order int) *model.Route {
	return &model.Route{
		ID:        id,
		Name:      id,
		Namespace: "default",
		Source:    model.SourceIngress,
		Title:     title,
		Group:     group,
		Order:     order,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestStoreSetAndGet(t *testing.T) {
	s := New(nil)
	r := newRoute("Ingress:default/app", "App", "default", 0)

	s.Set(r)

	got, ok := s.Get("Ingress:default/app")
	if !ok {
		t.Fatal("expected route to exist")
	}
	if got.Title != "App" {
		t.Errorf("got title %q, want %q", got.Title, "App")
	}
}

func TestStoreDelete(t *testing.T) {
	s := New(nil)
	r := newRoute("Ingress:default/app", "App", "default", 0)

	s.Set(r)
	s.Delete("Ingress:default/app")

	_, ok := s.Get("Ingress:default/app")
	if ok {
		t.Fatal("expected route to be deleted")
	}
}

func TestStoreListFiltersHidden(t *testing.T) {
	s := New(nil)
	s.Set(newRoute("r1", "Visible", "default", 0))

	hidden := newRoute("r2", "Hidden", "default", 0)
	hidden.Hidden = true
	s.Set(hidden)

	routes := s.List()
	if len(routes) != 1 {
		t.Fatalf("got %d routes, want 1", len(routes))
	}
	if routes[0].Title != "Visible" {
		t.Errorf("got title %q, want %q", routes[0].Title, "Visible")
	}
}

func TestStoreListSortOrder(t *testing.T) {
	s := New(nil)
	s.Set(newRoute("r1", "Zebra", "b-group", 0))
	s.Set(newRoute("r2", "Alpha", "a-group", 10))
	s.Set(newRoute("r3", "Beta", "a-group", 5))

	routes := s.List()
	if len(routes) != 3 {
		t.Fatalf("got %d routes, want 3", len(routes))
	}
	// a-group first (order 5 before 10), then b-group
	if routes[0].Title != "Beta" {
		t.Errorf("routes[0] = %q, want Beta", routes[0].Title)
	}
	if routes[1].Title != "Alpha" {
		t.Errorf("routes[1] = %q, want Alpha", routes[1].Title)
	}
	if routes[2].Title != "Zebra" {
		t.Errorf("routes[2] = %q, want Zebra", routes[2].Title)
	}
}

func TestStoreNotifyOnAdd(t *testing.T) {
	var events []ChangeEvent
	s := New(func(e ChangeEvent) { events = append(events, e) })

	s.Set(newRoute("r1", "App", "default", 0))

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
	if events[0].Type != ChangeAdded {
		t.Errorf("got type %q, want %q", events[0].Type, ChangeAdded)
	}
}

func TestStoreNotifyOnUpdate(t *testing.T) {
	var events []ChangeEvent
	s := New(func(e ChangeEvent) { events = append(events, e) })

	r := newRoute("r1", "App", "default", 0)
	s.Set(r)

	r2 := newRoute("r1", "App Updated", "default", 0)
	s.Set(r2)

	if len(events) != 2 {
		t.Fatalf("got %d events, want 2", len(events))
	}
	if events[1].Type != ChangeUpdated {
		t.Errorf("got type %q, want %q", events[1].Type, ChangeUpdated)
	}
}

func TestStoreNoNotifyOnSameData(t *testing.T) {
	var events []ChangeEvent
	s := New(func(e ChangeEvent) { events = append(events, e) })

	r := newRoute("r1", "App", "default", 0)
	s.Set(r)

	r2 := newRoute("r1", "App", "default", 0)
	s.Set(r2)

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1 (no update for unchanged data)", len(events))
	}
}

func TestStoreNotifyOnDelete(t *testing.T) {
	var events []ChangeEvent
	s := New(func(e ChangeEvent) { events = append(events, e) })

	s.Set(newRoute("r1", "App", "default", 0))
	s.Delete("r1")

	if len(events) != 2 {
		t.Fatalf("got %d events, want 2", len(events))
	}
	if events[1].Type != ChangeDeleted {
		t.Errorf("got type %q, want %q", events[1].Type, ChangeDeleted)
	}
}

func TestStoreNamespaces(t *testing.T) {
	s := New(nil)
	r1 := newRoute("r1", "A", "g", 0)
	r1.Namespace = "alpha"
	r2 := newRoute("r2", "B", "g", 0)
	r2.Namespace = "beta"
	r3 := newRoute("r3", "C", "g", 0)
	r3.Namespace = "alpha"

	s.Set(r1)
	s.Set(r2)
	s.Set(r3)

	ns := s.Namespaces()
	if len(ns) != 2 {
		t.Fatalf("got %d namespaces, want 2", len(ns))
	}
	if ns[0] != "alpha" || ns[1] != "beta" {
		t.Errorf("got namespaces %v, want [alpha beta]", ns)
	}
}

func TestStoreGroupedRoutes(t *testing.T) {
	s := New(nil)
	s.Set(newRoute("r1", "A", "infra", 0))
	s.Set(newRoute("r2", "B", "apps", 0))
	s.Set(newRoute("r3", "C", "infra", 0))

	grouped := s.GroupedRoutes()
	if len(grouped) != 2 {
		t.Fatalf("got %d groups, want 2", len(grouped))
	}
	if len(grouped["infra"]) != 2 {
		t.Errorf("infra has %d routes, want 2", len(grouped["infra"]))
	}
	if len(grouped["apps"]) != 1 {
		t.Errorf("apps has %d routes, want 1", len(grouped["apps"]))
	}
}
