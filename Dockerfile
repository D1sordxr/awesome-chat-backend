# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build all binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/outbox-processor ./cmd/outbox-processor/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/ws-server ./cmd/ws-server/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/topic-creator ./cmd/topic-creator/main.go

# Final lightweight image
FROM alpine:3.18

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api /app/api
COPY --from=builder /app/outbox-processor /app/outbox-processor
COPY --from=builder /app/ws-server /app/ws-server
COPY --from=builder /app/topic-creator /app/topic-creator

# Copy configs
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

# Create a non-root user and switch to it
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser