FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api /app/cmd/api/main.go

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Install necessary tools
RUN apk add --no-cache curl

# Copy the binary from the builder stage
COPY --from=builder /app/bin/api /app/api

# Copy the .env file
COPY .env /app/.env

# Expose the port
EXPOSE 3000

# Run the application
CMD ["/app/api"]
