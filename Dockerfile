# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/packify .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/packify .

# Copy template files and static assets
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./packify"]
