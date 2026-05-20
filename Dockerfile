# Stage 1: Builder
FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# Stage 2: Runtime (minimal)
FROM alpine:3.18

RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /root/

COPY --from=builder /app/main .

# Accept PORT_GRPC at build time; default 50051
ARG PORT_GRPC=50051
ENV PORT_GRPC=${PORT_GRPC}

EXPOSE ${PORT_GRPC}

HEALTHCHECK --interval=1m --timeout=10s --start-period=60s --retries=3 \
  CMD nc -z localhost ${PORT_GRPC} >/dev/null || exit 1

CMD ["./main"]
