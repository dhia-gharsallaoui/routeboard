package server

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/store"
)

type Server struct {
	cfg     *config.Config
	httpSrv *http.Server
}

func New(cfg *config.Config, s *store.Store, broker *SSEBroker, webFS fs.FS) *Server {
	handlers := NewHandlers(cfg, s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/routes", handlers.APIRoutes)
	mux.HandleFunc("GET /api/config", handlers.APIConfig)
	mux.HandleFunc("GET /health", handlers.Health)
	mux.Handle("GET /api/events", broker)

	if webFS != nil {
		mux.Handle("GET /", spaHandler(webFS))
	}

	return &Server{
		cfg: cfg,
		httpSrv: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 0, // disabled for SSE
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.httpSrv.Shutdown(shutdownCtx)
	}()

	slog.Info("starting HTTP server", "addr", s.httpSrv.Addr)
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// spaHandler serves the React SPA. It serves files from webFS if they exist,
// otherwise falls back to index.html for client-side routing.
func spaHandler(webFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(webFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else if path[0] == '/' {
			path = path[1:]
		}

		// Try to serve the file directly
		if _, err := fs.Stat(webFS, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fall back to index.html for SPA routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
