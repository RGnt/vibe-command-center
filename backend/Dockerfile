# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Disable CGO since we are using pure Go postgres driver (lib/pq)
ENV CGO_ENABLED=0

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Run the binary
FROM alpine:latest

WORKDIR /app

# Install dependencies required by the compiled binary (like libc)
RUN apk add --no-cache libc6-compat

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./main"]
