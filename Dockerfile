# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install git for go-git operations
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gist-sync ./cmd/syncd

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates git

# Create non-root user
RUN addgroup -g 1000 gist && \
    adduser -D -u 1000 -G gist gist

# Copy binary from builder
COPY --from=builder /build/gist-sync /app/

# Create work directory
RUN mkdir -p /tmp/gist-sync && chown gist:gist /tmp/gist-sync

# Switch to non-root user
USER gist

# Run the binary
ENTRYPOINT ["./gist-sync"]
