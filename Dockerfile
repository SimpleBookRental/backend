# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o book-rental-api ./cmd/api

# Final stage
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/book-rental-api .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy .env file for container
COPY .env ./.env

# Expose port from environment variable
EXPOSE ${SERVER_PORT}

# Run the application
ENTRYPOINT ["./book-rental-api"]
