# Build stage
FROM golang:alpine AS builder

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files from parent directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire project (we need models, config, etc.)
COPY . .

# Build the seeder
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seeder ./cmd/seed/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the seeder binary from builder stage
COPY --from=builder /app/seeder .

# Change ownership to appuser
RUN chown appuser:appuser /app/seeder

# Switch to non-root user
USER appuser

# Run the seeder
CMD ["./seeder"]