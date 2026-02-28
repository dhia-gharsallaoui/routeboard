# Stage 1: Build React frontend
FROM oven/bun:1 AS web-build
WORKDIR /web
COPY web/package.json web/bun.lock* ./
RUN bun install --frozen-lockfile
COPY web/ .
RUN bun run build

# Stage 2: Build Go binary
FROM golang:1.26-alpine AS go-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-build /web/dist ./internal/server/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /routeboard ./cmd/routeboard

# Stage 3: Runtime
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=go-build /routeboard /routeboard
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/routeboard"]
