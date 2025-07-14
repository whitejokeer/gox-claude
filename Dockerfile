# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o gox ./cmd/gox

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S gox && \
    adduser -u 1000 -S gox -G gox

# Copy binary from builder
COPY --from=builder /build/gox /usr/local/bin/gox

# Make binary executable
RUN chmod +x /usr/local/bin/gox

# Switch to non-root user
USER gox

# Set working directory
WORKDIR /app

# Expose default port
EXPOSE 3000

# Default command
ENTRYPOINT ["gox"]
CMD ["--help"]