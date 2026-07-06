# routeboard

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

Kubernetes-native Service Entry Portal — auto-discovers Ingress and HTTPRoute resources

**Homepage:** <https://github.com/dhia-gharsallaoui/routeboard>

## Installing

```bash
helm install routeboard oci://ghcr.io/dhia-gharsallaoui/helm/routeboard \
  -n routeboard --create-namespace
```

Or from this repository:

```bash
helm install routeboard deploy/helm/routeboard -n routeboard --create-namespace
```

## Exposing the UI

RouteBoard can be exposed through an `Ingress` or a Gateway API `HTTPRoute`.

```bash
# Ingress
helm install routeboard deploy/helm/routeboard \
  --set ingress.enabled=true --set ingress.hosts[0].host=routeboard.example.com

# Gateway API (requires the gateway.networking.k8s.io/v1 CRDs and a Gateway)
helm install routeboard deploy/helm/routeboard \
  --set gatewayAPI.enabled=true \
  --set gatewayAPI.httproute.parentRefs[0].name=my-gateway \
  --set gatewayAPI.httproute.parentRefs[0].namespace=gateway-system \
  --set gatewayAPI.httproute.hostnames[0]=routeboard.example.com
```

> **Set `parentRefs[].namespace` to where the Gateway lives.** It defaults to the
> release namespace, but Gateways usually run in a separate namespace (e.g.
> `gateway-system`, `istio-ingress`). If it points at the wrong namespace the route
> never attaches and the Gateway returns `404` (Envoy `NR` — no route). Verify with
> `kubectl get httproute -A` — the route should report `Accepted=True`.

`gatewayAPI.enabled` controls only the chart's own `HTTPRoute`. Discovery of other
HTTPRoutes is independent and on by default (`config.watchHTTPRoute`).

## Extra manifests

Deploy arbitrary resources alongside the chart via `extraManifests` — Gateway API
traffic policies, NetworkPolicies, or anything the chart doesn't template natively.
Each entry is rendered through `tpl`, so it can reference values and the release:

```yaml
extraManifests:
  - apiVersion: gateway.kgateway.dev/v1alpha1
    kind: TrafficPolicy
    metadata:
      name: '{{ include "routeboard.fullname" . }}'
    spec:
      targetRefs:
        - group: gateway.networking.k8s.io
          kind: HTTPRoute
          name: '{{ include "routeboard.fullname" . }}'
      retry:
        attempts: 3
```

## Testing

```bash
helm test routeboard -n routeboard   # in-cluster connectivity test
make helm-unittest                   # offline unit tests (helm-unittest plugin)
```

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| dhia-gharsallaoui |  | <https://github.com/dhia-gharsallaoui> |

## Source Code

* <https://github.com/dhia-gharsallaoui/routeboard>

## Requirements

Kubernetes: `>=1.21.0-0`

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity rules for pod scheduling. |
| config.labelSelector | string | `""` | Label selector to filter discovered routes (empty = all). |
| config.logLevel | string | `"info"` | Log level: debug, info, warn, or error. |
| config.namespaceAllowlist | string | `""` | Comma-separated namespaces to restrict discovery to (empty = all). |
| config.namespaceDenylist | string | `"kube-system,kube-public,kube-node-lease"` | Comma-separated namespaces to exclude from discovery. |
| config.port | int | `8080` | HTTP server port (container port). |
| config.resyncInterval | string | `"30m"` | Informer resync interval. |
| config.title | string | `"RouteBoard"` | Dashboard title. |
| config.watchHTTPRoute | bool | `true` | Watch HTTPRoute resources. Also gates the httproutes RBAC rule. |
| config.watchIngress | bool | `true` | Watch Ingress resources. Also gates the ingresses RBAC rule. |
| extraManifests | list | `[]` | Extra raw manifests to deploy alongside the chart. Each item may be a map or a string and is rendered through `tpl`, so it can reference values and the release (e.g. `{{ include "routeboard.fullname" . }}`). Useful for Gateway API traffic policies, NetworkPolicies, or any resource the chart doesn't template natively. |
| fullnameOverride | string | `""` | Override the fully-qualified release name used for resource names. |
| gatewayAPI.enabled | bool | `false` | Expose RouteBoard's own UI through a Gateway API HTTPRoute (parallel to ingress). This only controls the chart's own HTTPRoute — discovery of other HTTPRoutes is `config.watchHTTPRoute`. |
| gatewayAPI.httproute.annotations | object | `{}` | Annotations to add to the HTTPRoute. |
| gatewayAPI.httproute.hostnames | list | `[]` | Hostnames the route matches, e.g. `["routeboard.example.com"]`. |
| gatewayAPI.httproute.matches | list | `[{"path":{"type":"PathPrefix","value":"/"}}]` | Route match rules. Defaults to a PathPrefix match on `/`. |
| gatewayAPI.httproute.parentRefs | list | `[{"name":""}]` | Gateway(s) to attach to. At least one `parentRefs[].name` is required when `gatewayAPI.enabled` is true. `namespace` defaults to the release namespace. |
| gatewayAPI.verifyCRD | bool | `true` | Fail rendering if `gatewayAPI.enabled` but the Gateway API CRDs (`gateway.networking.k8s.io/v1`) are absent from the cluster. Set to false for offline `helm template`/CI (or pass `--api-versions gateway.networking.k8s.io/v1`). |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy. |
| image.repository | string | `"ghcr.io/dhia-gharsallaoui/routeboard"` | Container image repository. |
| image.tag | string | `""` | Image tag. Defaults to the chart `appVersion` when empty. |
| imagePullSecrets | list | `[]` | Secrets for pulling the image from a private registry. |
| ingress.annotations | object | `{}` | Annotations to add to the Ingress. |
| ingress.className | string | `""` | IngressClass name. |
| ingress.enabled | bool | `false` | Expose RouteBoard through a networking.k8s.io Ingress. |
| ingress.hosts | list | `[{"host":"routeboard.local","paths":[{"path":"/","pathType":"Prefix"}]}]` | Ingress host/path rules. |
| ingress.tls | list | `[]` | Ingress TLS configuration. |
| livenessProbe | object | `{"httpGet":{"path":"/health","port":"http"},"initialDelaySeconds":5,"periodSeconds":10}` | Liveness probe for the container. |
| nameOverride | string | `""` | Override the generated chart name portion of resource names. |
| nodeSelector | object | `{}` | Node selector for pod scheduling. |
| podAnnotations | object | `{}` | Extra annotations applied to the Pod template. |
| podDisruptionBudget.enabled | bool | `false` | Create a PodDisruptionBudget. Only effective with `replicaCount` > 1. |
| podDisruptionBudget.maxUnavailable | string | `""` | Maximum unavailable pods during a disruption. Set this OR `minAvailable`, not both. |
| podDisruptionBudget.minAvailable | int | `1` | Minimum available pods during a disruption. Mutually exclusive with `maxUnavailable`. |
| podLabels | object | `{}` | Extra labels applied to the Pod template. |
| podSecurityContext | object | `{"runAsNonRoot":true,"runAsUser":65534,"seccompProfile":{"type":"RuntimeDefault"}}` | Pod-level security context. `readOnlyRootFilesystem` is intentionally NOT here — it is a container-only field (see `securityContext`). |
| priorityClassName | string | `""` | PriorityClass name for the pod. |
| readinessProbe | object | `{"httpGet":{"path":"/health","port":"http"},"initialDelaySeconds":3,"periodSeconds":5}` | Readiness probe for the container. |
| replicaCount | int | `1` | Number of RouteBoard replicas to run. |
| resources | object | `{"limits":{"cpu":"100m","memory":"128Mi"},"requests":{"cpu":"50m","memory":"64Mi"}}` | Container resource requests and limits. |
| revisionHistoryLimit | int | `3` | Number of old ReplicaSets to retain for rollback. |
| securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true}` | Container-level security context. |
| service.annotations | object | `{}` | Annotations to add to the Service. |
| service.port | int | `80` | Service port. |
| service.type | string | `"ClusterIP"` | Service type. |
| serviceAccount.annotations | object | `{}` | Annotations to add to the ServiceAccount. |
| serviceAccount.automountServiceAccountToken | bool | `true` | Mount the API token. RouteBoard calls the Kubernetes API, so this must stay true. |
| serviceAccount.create | bool | `true` | Create a ServiceAccount for RouteBoard. |
| serviceAccount.name | string | `""` | Name of the ServiceAccount to use. Generated from the fullname when empty. |
| tmpDir.sizeLimit | string | `"16Mi"` | Size limit for the writable `/tmp` emptyDir (required because the root FS is read-only). Bounded so it cannot exhaust node ephemeral storage. |
| tolerations | list | `[]` | Tolerations for pod scheduling. |
| topologySpreadConstraints | list | `[]` | Topology spread constraints for pod scheduling (multi-replica zone/node spreading). |
