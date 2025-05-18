# Start from the official Go image as the build stage
FROM golang:1.21 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go app
RUN go build -o coupon-system main.go

# Start a new stage from a minimal image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/coupon-system .

# Expose port (adjust based on what your app listens on)
EXPOSE 8080

# Command to run the executable
CMD ["./coupon-system"]
