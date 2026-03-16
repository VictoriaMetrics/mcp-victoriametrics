# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Install Git for potential dependencies and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy source
COPY go.mod go.sum ./
COPY . .

# Download dependencies and build
RUN go mod download
RUN go build -o mcp-victoriametrics cmd/mcp-victoriametrics/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/mcp-victoriametrics /usr/local/bin/mcp-victoriametrics

# Default entrypoint
ENTRYPOINT ["mcp-victoriametrics"]
