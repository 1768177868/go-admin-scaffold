# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache tzdata

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# Set timezone
ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD ["./main"] 