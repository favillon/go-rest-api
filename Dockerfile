# Stage 1: Builder
FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# Stage 2: Runtime (minimal)
FROM alpine:3.18

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/main .

# Accept PORT_APP at build time; default 8082
ARG PORT_APP=8082
ENV PORT_APP=${PORT_APP}

EXPOSE ${PORT_APP}

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:${PORT_APP}/api/v1/graphql || exit 1

CMD ["./main"]
