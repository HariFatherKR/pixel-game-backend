# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Copy source code first to run go mod tidy
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy

# Install swag and generate Swagger docs (allow failure)
RUN go install github.com/swaggo/swag/cmd/swag@latest || true
RUN swag init -g cmd/server/main.go --parseDependency --parseInternal -o ./docs --ot json,yaml || true

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy swagger docs (copy entire directory to ensure generated files are included)
COPY --from=builder /app/docs ./docs

# Copy config file if needed
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]