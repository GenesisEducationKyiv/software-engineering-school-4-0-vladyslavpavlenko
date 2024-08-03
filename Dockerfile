# Build stage
FROM golang:1.22.2 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o apiApp ./cmd/api

# Final stage
FROM alpine:latest

# Create necessary directories
RUN mkdir -p /var/log/app

# Copy the built binary and other necessary files
COPY --from=builder /app/apiApp /app/apiApp
COPY --from=builder /app/.env /app/.env

WORKDIR /app

# Run the application
CMD ["/app/apiApp"]