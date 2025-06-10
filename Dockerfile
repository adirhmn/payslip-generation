# -------------------
# Build stage
# -------------------
FROM golang:1.23 as builder

# Enable Go modules
ENV GO111MODULE=on

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application binary
# Use CGO_ENABLED=0 to statically link (optional but good for portability)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app/main.go

# -------------------
# Final runtime stage
# -------------------
FROM debian:bullseye-slim

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/app .

# Optionally copy .env or config files if used at runtime
COPY .env .

# Run the application
CMD ["./app"]
    