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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# Stage 2: Minimal runtime image with Nginx
FROM nginx:alpine

# Install necessary packages
RUN apk --no-cache add ca-certificates bash

# Set the working directory
WORKDIR /app

# Copy the Nginx configuration file into the container
COPY nginx/conf/nginx.conf /etc/nginx/nginx.conf

# Copy the SSL certificates
COPY ./nginx/certificate /etc/nginx/certificate

# Copy the Go app binary from the builder stage
COPY --from=builder /app/app .

# Expose ports for Nginx and the Go backend
EXPOSE 80
EXPOSE 443
EXPOSE 5000

# Start the Go application in the background and Nginx in the foreground
CMD ./app & nginx -g "daemon off;"
