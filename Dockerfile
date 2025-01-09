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

# Create a non-root user and set ownership
RUN adduser -D appuser && \
    chown -R appuser:appuser /app

# Make the binary executable
RUN chmod +x /app/app

USER appuser

# Default environment variables
ENV PORT=8080
ENV DB_URI=""

# Expose the port
EXPOSE ${PORT}

# Directly run the application
CMD ["./app"]