FROM golang:1.23.4-alpine AS builder

WORKDIR /build

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire project
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o client ./client

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/client .

# Set permissions
RUN chmod +x /app/client && \
    chown -R nobody:nobody /app

# Switch to non-root user
USER nobody

EXPOSE 8081

CMD ["./client"]