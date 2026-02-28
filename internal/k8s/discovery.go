package k8s

import (
	"context"
	"log/slog"

	"k8s.io/client-go/kubernetes"
)

func gatewayAPIAvailable(ctx context.Context, client kubernetes.Interface) bool {
	_, err := client.Discovery().ServerResourcesForGroupVersion("gateway.networking.k8s.io/v1")
	if err != nil {
		slog.Warn("Gateway API CRDs not found, HTTPRoute watching disabled", "error", err)
		return false
	}
	return true
}
