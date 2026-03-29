# RouteBoard — Feature Inventory

> Living document tracking implemented features and future ideas.
> Contributions welcome — open an issue to discuss any of these!

---

## Legend

- [x] Implemented
- [ ] Planned / idea

---

## Core Discovery

- [x] `networking.k8s.io/v1 Ingress` auto-discovery
- [x] `gateway.networking.k8s.io/v1 HTTPRoute` auto-discovery
- [x] Real-time updates via SharedInformerFactory (no polling)
- [x] Namespace allow/deny list filtering
- [x] Label selector filtering
- [x] Routes without URLs filtered from display
- [x] Gateway API CRD graceful degradation (works if CRDs not installed)
- [ ] `IngressRoute` (Traefik CRD) support
- [ ] `VirtualService` (Istio) support
- [ ] `HTTPProxy` (Contour) support
- [ ] Generic CRD support via dynamic informers + config
- [ ] Multi-cluster aggregation (multiple kubeconfigs, cluster badge on cards)
- [ ] Service discovery mode (watch `v1/Service` with annotations, no Ingress required)

## Annotations

- [x] `routeboard.io/title` — custom display name
- [x] `routeboard.io/description` — short description
- [x] `routeboard.io/icon` — override auto-detected icon
- [x] `routeboard.io/group` — custom grouping (default: namespace)
- [x] `routeboard.io/order` — sort order within group
- [x] `routeboard.io/hidden` — hide from dashboard
- [x] `routeboard.io/url` — override computed URL
- [ ] `routeboard.io/health-path` — custom health check endpoint (default: `/`)
- [ ] `routeboard.io/health-method` — GET vs HEAD (default: HEAD)
- [ ] `routeboard.io/health-expected-status` — expected HTTP status code
- [ ] `routeboard.io/tags` — comma-separated tags for filtering
- [ ] `routeboard.io/category` — service category (monitoring, storage, etc.)
- [ ] `routeboard.io/docs` — link to documentation
- [ ] `routeboard.io/owner` — team/person responsible
- [ ] Custom annotation prefix via config (not just `routeboard.io/`)

## UI / UX

### Views & Layout

- [x] Grid view (responsive: 1–4 columns)
- [x] List view (compact rows)
- [x] Dark mode (default) with light mode toggle
- [x] Grouped by namespace or custom annotation
- [x] Group headers with route count
- [x] Loading skeleton with shimmer animation
- [x] Empty state with helpful kubectl command
- [x] Mobile responsive layout
- [x] Sticky header with backdrop blur
- [ ] Compact/dense view mode (more cards, less padding)
- [ ] Group collapse/expand (remember state in localStorage)
- [ ] Dashboard stats bar (total routes, healthy %, avg response time)
- [ ] Fullscreen mode (hide header, maximize cards)
- [ ] Sort options (name, health, last updated, response time)
- [ ] Drag-and-drop reordering (persisted to localStorage)

### Search & Filtering

- [x] Instant client-side search (title, URL, description, name, namespace)
- [x] Namespace dropdown filter
- [x] Health status filter (healthy / degraded / unhealthy)
- [ ] Tag/label filter chips
- [ ] Combined filter URL params (shareable filter state)
- [ ] Saved filter presets
- [ ] "Show hidden routes" toggle for admins

### Cards & Interaction

- [x] Brand SVG icons from Simple Icons CDN with fallback
- [x] Auto-guess icons for unknown services
- [x] Brand colors from CDN
- [x] TLS lock indicator
- [x] Source badge (Ingress / HTTPRoute)
- [x] Namespace badge
- [x] Click-to-copy URL with checkmark feedback
- [x] External link arrow on hover
- [x] Card hover animation (golden glow, lift)
- [x] Staggered entrance animations
- [x] Favorites / pin routes to top (localStorage)
- [x] Health status dots (green/yellow/red/gray) with tooltip
- [x] Uptime sparkline (last 60 checks)
- [ ] Route detail panel (click to expand: all metadata, full health history, response times)
- [ ] QR code for mobile access to any route
- [ ] "Open in new tab" keyboard shortcut
- [ ] Notes/comments on routes (localStorage)
- [ ] Recently visited routes section

### Navigation & Power User

- [ ] Keyboard navigation (`/` to search, `j`/`k` to navigate, `Enter` to open)
- [ ] Command palette (`Cmd+K` / `Ctrl+K`) — search, filter, navigate, toggle theme
- [ ] Vim-style key bindings option
- [ ] Browser tab title shows health summary (e.g. "RouteBoard — 2 unhealthy")
- [ ] URL deep linking to filtered views (`?namespace=monitoring&health=unhealthy`)

### Theming & Branding

- [ ] Custom logo (via Helm values or annotation)
- [ ] Custom CSS injection
- [ ] Custom accent color
- [ ] Company name / subtitle
- [ ] Favicon override
- [ ] Multiple built-in themes beyond dark/light

## Health Monitoring

- [x] Periodic HEAD requests to each route URL
- [x] Configurable interval and timeout
- [x] TLS verification skipped (for internal certs)
- [x] Healthy / degraded / unhealthy / unknown states
- [x] Health status dots with time-ago tooltip
- [x] Health history ring buffer (last 60 checks)
- [x] Uptime sparkline visualization
- [x] Health filter in header
- [ ] Custom health check path per route (annotation)
- [ ] Custom health check method (GET/HEAD) per route
- [ ] Expected status code override per route
- [ ] Response time tracking and display (ms on card)
- [ ] Response time sparkline
- [ ] SSL certificate expiry monitoring and warning
- [ ] Health check results persistence (optional SQLite/file)
- [ ] Uptime percentage display (24h / 7d / 30d)
- [ ] Incident timeline (when did each transition happen)
- [ ] Scheduled maintenance mode (suppress alerts, show banner)
- [ ] Health check dependencies (mark a route as dependent on another)

## Notifications

- [x] Webhook on health state transitions
- [x] JSON payload format
- [x] Slack Block Kit format
- [x] Discord embed format
- [ ] Microsoft Teams webhook format
- [ ] Email notifications (SMTP)
- [ ] PagerDuty integration
- [ ] Opsgenie integration
- [ ] Ntfy.sh support
- [ ] Gotify support
- [ ] Notification cooldown / debounce (don't spam on flapping)
- [ ] Notification rules (only specific routes or namespaces)
- [ ] Recovery notifications ("X is back up after 5m downtime")
- [ ] Daily/weekly health digest
- [ ] Notification history log (in-memory, viewable in UI)

## API & Integration

- [x] `GET /api/routes` — JSON route list with filtering
- [x] `GET /api/config` — dashboard config
- [x] `GET /api/events` — SSE stream
- [x] `GET /health` — liveness probe
- [ ] `GET /api/routes/:id` — single route detail
- [ ] `GET /api/routes/:id/health` — health history for a route
- [ ] `GET /metrics` — Prometheus metrics endpoint
  - `routeboard_routes_total` (gauge, by namespace/source)
  - `routeboard_routes_healthy` (gauge)
  - `routeboard_routes_unhealthy` (gauge)
  - `routeboard_health_check_duration_seconds` (histogram)
  - `routeboard_sse_clients_connected` (gauge)
- [ ] OpenAPI / Swagger spec
- [ ] JSON export of entire catalog
- [ ] Prometheus `ServiceMonitor` in Helm chart
- [ ] Grafana dashboard JSON template
- [ ] Backstage catalog integration plugin
- [ ] Browser extension (new tab page showing RouteBoard)
- [ ] Embeddable widget mode (iframe with minimal chrome)
- [ ] Terraform data source (read routes programmatically)

## Security & Access Control

- [ ] OIDC authentication (Keycloak, Authentik, Dex)
- [ ] OAuth2 proxy compatibility (trust `X-Forwarded-User` header)
- [ ] Basic auth (username/password from secret)
- [ ] API key authentication for `/api/*` endpoints
- [ ] Role-based visibility (only show routes in namespaces user has access to)
- [ ] Read-only mode (disable health checks, webhooks)
- [ ] Audit log (who accessed what, when)
- [ ] Content Security Policy headers
- [ ] Rate limiting on API endpoints

## Deployment & Operations

- [x] Helm chart with RBAC, ServiceAccount, security context
- [x] Multi-stage Dockerfile (bun + Go + distroless)
- [x] Multi-arch images (amd64 + arm64)
- [x] GitHub Actions CI (Go lint/test, Biome lint/typecheck)
- [x] GitHub Actions release (Docker + Helm OCI to GHCR)
- [x] `go:embed` single binary with frontend
- [x] Graceful shutdown
- [ ] ArtifactHub listing for Helm chart
- [ ] Kustomize manifests (alternative to Helm)
- [ ] Operator with CRD (`RouteBoard` custom resource for multi-instance)
- [ ] HA mode (multiple replicas, leader election for health checks)
- [ ] Pod disruption budget in Helm chart
- [ ] Network policy in Helm chart
- [ ] Resource auto-tuning (adjust limits based on route count)
- [ ] Liveness + readiness + startup probes (startup currently missing)
- [ ] ConfigMap-based config (alternative to env vars)
- [ ] HPA (horizontal pod autoscaler) template

## Developer Experience

- [ ] `make dev` — concurrent Go + Vite with hot reload
- [ ] Mock data mode (no cluster needed, fake routes for UI development)
- [ ] E2E tests with kind cluster in CI
- [ ] Storybook for UI components
- [ ] Contributing guide (`CONTRIBUTING.md`)
- [ ] Architecture decision records (`docs/adr/`)
- [ ] Development container (devcontainer.json)
- [ ] Pre-commit hooks (lint, format, test)

## Creative / Unique Ideas

These are more ambitious features that would make RouteBoard stand out:

- [ ] **Service topology map** — visual graph showing services and their connections (inferred from Gateway parentRefs, or via annotations)
- [ ] **Changelog feed** — track when routes were added, removed, or changed. Show a timeline view. "Grafana was added 3 days ago by namespace monitoring"
- [ ] **Resource usage overlay** — show CPU/memory of the backing pod next to the route (requires additional RBAC for pods/metrics)
- [ ] **Status page mode** — public-facing `/status` endpoint showing uptime for selected services, embeddable in external status pages
- [ ] **Bookmarklet generator** — one click to generate a browser bookmark bar with all your routes organized by group
- [ ] **Traffic light widget** — single aggregate indicator (all green / some yellow / any red) for embedding in Slack status, Notion, etc.
- [ ] **"What's down?" AI summary** — natural language summary of current cluster health ("3 services are degraded in the monitoring namespace, likely due to the Prometheus upgrade 2h ago")
- [ ] **Comparison mode** — diff two clusters side by side, highlight routes that exist in one but not the other
- [ ] **Route dependencies** — annotate dependencies between services, show impact radius when something goes down
- [ ] **Maintenance windows** — schedule planned downtime for a route, auto-suppress health alerts, show a "maintenance" badge
- [ ] **Daily morning digest** — email/Slack summary at 9am: "Your cluster has 15 routes, 14 healthy, Grafana has been down for 2h"
- [ ] **PWA support** — installable as a home screen app on mobile, with offline shell and push notifications for health alerts
