package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/store"
)

type Handlers struct {
	cfg   *config.Config
	store *store.Store
}

func NewHandlers(cfg *config.Config, store *store.Store) *Handlers {
	return &Handlers{cfg: cfg, store: store}
}

type configResponse struct {
	Title      string   `json:"title"`
	Namespaces []string `json:"namespaces"`
}

func (h *Handlers) APIRoutes(w http.ResponseWriter, r *http.Request) {
	routes := h.store.List()

	if ns := r.URL.Query().Get("namespace"); ns != "" {
		filtered := routes[:0]
		for _, route := range routes {
			if route.Namespace == ns {
				filtered = append(filtered, route)
			}
		}
		routes = filtered
	}

	if q := strings.ToLower(r.URL.Query().Get("q")); q != "" {
		filtered := routes[:0]
		for _, route := range routes {
			if strings.Contains(strings.ToLower(route.Title), q) ||
				strings.Contains(strings.ToLower(route.URL), q) ||
				strings.Contains(strings.ToLower(route.Description), q) ||
				strings.Contains(strings.ToLower(route.Name), q) {
				filtered = append(filtered, route)
			}
		}
		routes = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(routes)
}

func (h *Handlers) APIConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(configResponse{
		Title:      h.cfg.Title,
		Namespaces: h.store.Namespaces(),
	})
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
