# syntax=docker/dockerfile:1

# Stage 1: Build the Go binary
FROM golang:1.23.5-alpine AS builder
WORKDIR /app

# Copy all Go source files and go.mod/go.sum files
COPY . .

# Download dependencies and build the binary
RUN go mod tidy
RUN go build -o chitest

# Stage 2: Create a minimal image with the compiled binary
FROM alpine:latest
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/chitest .

# Expose a port if your application uses one (optional)
EXPOSE 8080

# Set the entrypoint to run your binary
ENTRYPOINT ["./chitest"]