# Stage 1: Build Go Application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the entire application source code
COPY . .

# Build the application binary (static build for portability)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

# Stage 2: Create a lightweight image for running the application
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates bash

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Expose the application port
EXPOSE 5000

# Command to run the application
CMD ["./app"]
