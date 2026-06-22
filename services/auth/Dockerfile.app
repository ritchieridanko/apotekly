# ---------- Build Stage ----------
FROM golang:1.25.0-alpine3.22 AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set work directory
WORKDIR /app/auth

# Copy and download app dependencies
COPY services/shared/go.mod services/shared/go.sum ../shared/
COPY services/auth/go.mod services/auth/go.sum ./
RUN go mod download

# Copy app source
COPY services/shared/configs ../shared/configs
COPY services/shared/constants ../shared/constants
COPY services/shared/contract/apis/v1 ../shared/contract/apis/v1
COPY services/shared/contract/events/v1 ../shared/contract/events/v1
COPY services/shared/infra/cache ../shared/infra/cache
COPY services/shared/infra/database ../shared/infra/database
COPY services/shared/infra/logger ../shared/infra/logger
COPY services/shared/infra/publisher ../shared/infra/publisher
COPY services/shared/infra/tracer ../shared/infra/tracer
COPY services/shared/utils ../shared/utils
COPY services/auth/cmd/app ./cmd/app
COPY services/auth/configs ./configs
COPY services/auth/internal ./internal

# Build binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/app/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.22

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set work directory
WORKDIR /root

# Copy from the Build Stage
COPY --from=builder /app/auth/bin ./bin
COPY --from=builder /app/auth/configs ./configs

# Expose port
EXPOSE 50051

# Set entry point
ENTRYPOINT ["./bin/app"]
