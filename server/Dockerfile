# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./server

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache curl

WORKDIR /root/

COPY --from=builder /app/server .

ENV REDIS_URL="192.168.1.6:6379"
ENV PORT 8081

CMD ["./server"]