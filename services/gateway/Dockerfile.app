# ---------- Build Stage ----------
FROM golang:1.25.0-alpine3.22 AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set work directory
WORKDIR /app/gateway

# Copy and download app dependencies
COPY services/shared/go.mod services/shared/go.sum ../shared/
COPY services/gateway/go.mod services/gateway/go.sum ./
RUN go mod download

# Copy app source
COPY services/shared/configs ../shared/configs
COPY services/shared/constants ../shared/constants
COPY services/shared/contract/apis/v1 ../shared/contract/apis/v1
COPY services/shared/infra/logger ../shared/infra/logger
COPY services/shared/infra/services ../shared/infra/services
COPY services/shared/infra/tracer ../shared/infra/tracer
COPY services/shared/utils ../shared/utils
COPY services/gateway/cmd/app ./cmd/app
COPY services/gateway/configs ./configs
COPY services/gateway/internal ./internal

# Build binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.22

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set work directory
WORKDIR /root

# Copy from the Build Stage
COPY --from=builder /app/gateway/bin ./bin
COPY --from=builder /app/gateway/configs ./configs

# Expose port
EXPOSE 8080

# Set entry point
ENTRYPOINT ["./bin/app"]
