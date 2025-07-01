# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (needed if you have dependencies)
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source files
COPY . .

# Build the binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

# Stage 2: Run
FROM scratch

# Copy binary and necessary files from builder
COPY --from=builder /app/app /app/app

# Set working directory
WORKDIR /app

# Expose port (adjust if needed)
EXPOSE 5001

# Command to run
CMD ["./app"]
