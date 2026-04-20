# Stage 1: Builder
FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# Stage 2: Runtime (minimal)
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the compiled binary from builder
COPY --from=builder /app/main .

# Runtime port used by app and healthcheck
ENV PORT_APP=8082

# Expose the port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:${PORT_APP:-8082}/api/v1/productos || exit 1

# Run the application
CMD ["./main"]
