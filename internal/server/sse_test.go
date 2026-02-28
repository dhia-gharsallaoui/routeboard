package server

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dhia/routeboard/internal/model"
	"github.com/dhia/routeboard/internal/store"
)

func TestSSEBrokerBroadcast(t *testing.T) {
	broker := NewSSEBroker()

	// Start SSE server
	srv := httptest.NewServer(broker)
	defer srv.Close()

	// Connect client
	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", ct)
	}

	reader := bufio.NewReader(resp.Body)

	// Read the initial "connected" event
	readSSEEvent(t, reader, "connected")

	// Send a route change event
	broker.Notify(store.ChangeEvent{
		Type:  store.ChangeAdded,
		Route: &model.Route{ID: "test-1", Title: "Test Route"},
	})

	// Read the route-change event
	readSSEEvent(t, reader, "route-change")
}

func TestSSEBrokerMultipleClients(t *testing.T) {
	broker := NewSSEBroker()
	srv := httptest.NewServer(broker)
	defer srv.Close()

	// Connect two clients
	resp1, _ := http.Get(srv.URL)
	defer resp1.Body.Close()
	resp2, _ := http.Get(srv.URL)
	defer resp2.Body.Close()

	reader1 := bufio.NewReader(resp1.Body)
	reader2 := bufio.NewReader(resp2.Body)

	// Skip connected events
	readSSEEvent(t, reader1, "connected")
	readSSEEvent(t, reader2, "connected")

	// Broadcast
	broker.Notify(store.ChangeEvent{
		Type:  store.ChangeUpdated,
		Route: &model.Route{ID: "test-1", Title: "Updated"},
	})

	// Both should receive
	readSSEEvent(t, reader1, "route-change")
	readSSEEvent(t, reader2, "route-change")
}

func readSSEEvent(t *testing.T, reader *bufio.Reader, expectedEvent string) {
	t.Helper()

	done := make(chan struct{})
	var lines []string

	go func() {
		defer close(done)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimRight(line, "\r\n")
			if line == "" {
				return // end of event
			}
			lines = append(lines, line)
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for SSE event %q", expectedEvent)
	}

	if len(lines) == 0 {
		t.Fatalf("got no SSE lines, expected event %q", expectedEvent)
	}

	found := false
	for _, line := range lines {
		if strings.HasPrefix(line, "event: "+expectedEvent) {
			found = true
		}
	}
	if !found {
		t.Errorf("expected event %q in lines %v", expectedEvent, lines)
	}
}
