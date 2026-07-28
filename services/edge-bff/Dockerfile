# services/Edge-bff/Dockerfile
FROM golang:1.24.11-alpine AS builder
WORKDIR /app
# Install dependencies
RUN apk add --no-cache git gcc musl-dev
# Copy go mod files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy source code
COPY . .
# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go
# Production stage - Distroless
FROM gcr.io/distroless/static-debian11
WORKDIR /app
# Copy binary from builder
COPY --from=builder /app/main .
# Copy config files if needed
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./main"]
