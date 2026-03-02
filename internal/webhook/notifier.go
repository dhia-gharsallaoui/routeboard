package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dhia/routeboard/internal/model"
	"github.com/dhia/routeboard/internal/store"
)

type Notifier struct {
	url    string
	format string
	client *http.Client
}

func NewNotifier(url, format string) *Notifier {
	return &Notifier{
		url:    url,
		format: format,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type jsonPayload struct {
	Route    string `json:"route"`
	URL      string `json:"url"`
	Status   string `json:"status"`
	Previous string `json:"previous"`
	Time     string `json:"timestamp"`
}

// HandleEvent is called on every store change event. It only sends
// webhooks for health state transitions (e.g. healthy → unhealthy).
func (n *Notifier) HandleEvent(event store.ChangeEvent) {
	if event.Type != store.ChangeHealth {
		return
	}
	if event.PreviousHealth == "" || event.PreviousHealth == event.Route.Health {
		return
	}

	go n.send(event.Route, event.PreviousHealth)
}

func (n *Notifier) send(route *model.Route, previous model.HealthStatus) {
	var body []byte
	var contentType string
	var err error

	switch n.format {
	case "slack":
		body, err = n.slackPayload(route, previous)
		contentType = "application/json"
	case "discord":
		body, err = n.discordPayload(route, previous)
		contentType = "application/json"
	default:
		body, err = n.jsonPayload(route, previous)
		contentType = "application/json"
	}

	if err != nil {
		slog.Error("webhook: failed to marshal payload", "error", err)
		return
	}

	resp, err := n.client.Post(n.url, contentType, bytes.NewReader(body))
	if err != nil {
		slog.Error("webhook: failed to send", "error", err, "route", route.Title)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		slog.Warn("webhook: non-2xx response", "status", resp.StatusCode, "route", route.Title)
	} else {
		slog.Info("webhook: sent", "route", route.Title, "status", route.Health, "previous", previous)
	}
}

func (n *Notifier) jsonPayload(route *model.Route, previous model.HealthStatus) ([]byte, error) {
	return json.Marshal(jsonPayload{
		Route:    route.Title,
		URL:      route.URL,
		Status:   string(route.Health),
		Previous: string(previous),
		Time:     time.Now().UTC().Format(time.RFC3339),
	})
}

func (n *Notifier) slackPayload(route *model.Route, previous model.HealthStatus) ([]byte, error) {
	color := "#36a64f" // green
	switch route.Health {
	case model.HealthUnhealthy:
		color = "#dc2626"
	case model.HealthDegraded:
		color = "#f59e0b"
	}

	payload := map[string]any{
		"attachments": []map[string]any{
			{
				"color": color,
				"blocks": []map[string]any{
					{
						"type": "section",
						"text": map[string]string{
							"type": "mrkdwn",
							"text": fmt.Sprintf("*%s* is now *%s*\n(was %s)\n<%s|Open service>",
								route.Title, route.Health, previous, route.URL),
						},
					},
				},
			},
		},
	}
	return json.Marshal(payload)
}

func (n *Notifier) discordPayload(route *model.Route, previous model.HealthStatus) ([]byte, error) {
	color := 0x36a64f // green
	switch route.Health {
	case model.HealthUnhealthy:
		color = 0xdc2626
	case model.HealthDegraded:
		color = 0xf59e0b
	}

	payload := map[string]any{
		"embeds": []map[string]any{
			{
				"title":       fmt.Sprintf("%s is now %s", route.Title, route.Health),
				"description": fmt.Sprintf("Previous status: %s\nURL: %s", previous, route.URL),
				"color":       color,
				"timestamp":   time.Now().UTC().Format(time.RFC3339),
			},
		},
	}
	return json.Marshal(payload)
}
