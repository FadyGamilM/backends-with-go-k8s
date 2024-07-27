# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o poller ./poller

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache curl

WORKDIR /root/

COPY --from=builder /app/poller .

ENV SERVER_URL="http://192.168.1.6:8081/counter"

CMD ["./poller"]