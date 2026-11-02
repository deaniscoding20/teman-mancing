# Stage 1: Build binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Download dependencies first for layer caching
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api ./cmd/api

# Stage 2: Runtime image
FROM alpine:3.20

WORKDIR /app

# Install runtime utilities & certificates
RUN apk add --no-cache ca-certificates tzdata wget \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy application binary and migrations
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/migrations /app/migrations

# Secure permissions
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

ENTRYPOINT ["/app/api"]
