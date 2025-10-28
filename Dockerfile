# Build stage
FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -v -o /app/main ./cmd

# Runtime stage
FROM alpine:3.22

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main /app/main

# Run the application
CMD ["/app/main"]