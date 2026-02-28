package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dhia/routeboard/internal/config"
	"github.com/dhia/routeboard/internal/k8s"
	"github.com/dhia/routeboard/internal/server"
	"github.com/dhia/routeboard/internal/store"
)

var version = "dev"

func main() {
	cfg := config.Load()
	setupLogging(cfg.LogLevel)

	slog.Info("starting routeboard", "version", version)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	broker := server.NewSSEBroker()
	routeStore := store.New(broker.Notify)

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

	srv := server.New(cfg, routeStore, broker, server.WebFS())
	if err := srv.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}

	slog.Info("routeboard stopped")
}

func setupLogging(level string) {
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})))
}
