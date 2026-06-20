package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/health"
	"github.com/dhia/routeboard/internal/k8s"
	"github.com/dhia/routeboard/internal/server"
	"github.com/dhia/routeboard/internal/store"
	"github.com/dhia/routeboard/internal/webhook"
)

var version = "dev"

func main() {
	cfg := config.Load()
	setupLogging(cfg.LogLevel, cfg.LogFormat)

	slog.Info("starting routeboard", "version", version)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	broker := server.NewSSEBroker()
	listeners := []store.NotifyFunc{broker.Notify}

	if cfg.WebhookURL != "" {
		notifier := webhook.NewNotifier(cfg.WebhookURL, cfg.WebhookFormat)
		listeners = append(listeners, notifier.HandleEvent)
		slog.Info("webhook notifications enabled", "url", cfg.WebhookURL, "format", cfg.WebhookFormat)
	}

	routeStore := store.New(listeners...)

	clients, err := k8s.NewClients(cfg.Kubeconfig)
	if err != nil {
		slog.Error("failed to create kubernetes clients", "error", err)
		os.Exit(1)
	}

	watcher := k8s.NewWatcher(cfg, clients, routeStore)

	go func() {
		if err := watcher.Run(ctx); err != nil {
			slog.Error("watcher error", "error", err)
			cancel()
		}
	}()

	if cfg.HealthEnabled {
		checker := health.NewChecker(cfg, routeStore)
		go checker.Run(ctx)
	}

	srv := server.New(cfg, routeStore, broker, server.WebFS())
	if err := srv.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}

	slog.Info("routeboard stopped")
}

func setupLogging(level, format string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: logLevel}
	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(handler))
}
