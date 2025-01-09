# Stage 1: Build Go Application
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the entire application source code
COPY . .

# Build the application binary (static build for portability)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates bash

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Script to check environment variables and start the application
COPY <<'EOF' /app/start.sh
#!/bin/bash
if [ -z "${DB_URI}" ]; then
    echo "ERROR: DB_URI environment variable is required"
    exit 1
fi
exec ./app
EOF

# Set permissions before switching user
RUN chmod +x /app/start.sh && \
    chmod +x /app/app

# Create a non-root user and set ownership
RUN adduser -D appuser && \
    chown -R appuser:appuser /app
USER appuser

# Default environment variables with ability to override
ENV PORT=8080 \
    DB_URI=""

# Expose the port
EXPOSE ${PORT}

# Add a healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Use the start script as the entrypoint
ENTRYPOINT ["/app/start.sh"]