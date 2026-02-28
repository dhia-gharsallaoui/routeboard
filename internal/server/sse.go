package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/dhia/routeboard/internal/store"
)

type SSEBroker struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

func NewSSEBroker() *SSEBroker {
	return &SSEBroker{
		clients: make(map[chan []byte]struct{}),
	}
}

func (b *SSEBroker) Notify(event store.ChangeEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal SSE event", "error", err)
		return
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	for clientCh := range b.clients {
		select {
		case clientCh <- data:
		default:
			slog.Warn("dropping SSE event for slow client")
		}
	}
}

func (b *SSEBroker) subscribe() chan []byte {
	ch := make(chan []byte, 16)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	slog.Debug("SSE client connected", "total", len(b.clients))
	return ch
}

func (b *SSEBroker) unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
	close(ch)
	slog.Debug("SSE client disconnected", "total", len(b.clients))
}

func (b *SSEBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	clientCh := b.subscribe()
	defer b.unsubscribe(clientCh)

	fmt.Fprintf(w, "event: connected\ndata: {}\n\n")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-clientCh:
			if !ok {
				return
			}
			fmt.Fprintf(w, "event: route-change\ndata: %s\n\n", data)
			flusher.Flush()
		}
	}
}
