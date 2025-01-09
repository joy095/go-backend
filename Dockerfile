# Stage 1: Build Go Application
FROM golang:1.22-alpine AS builder

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

# Create start script
RUN echo '#!/bin/bash' > /app/start.sh && \
    echo 'if [ -z "${DB_URI}" ]; then' >> /app/start.sh && \
    echo '    echo "ERROR: DB_URI environment variable is required"' >> /app/start.sh && \
    echo '    exit 1' >> /app/start.sh && \
    echo 'fi' >> /app/start.sh && \
    echo 'exec ./app' >> /app/start.sh

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