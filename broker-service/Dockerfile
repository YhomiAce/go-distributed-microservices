# Build
FROM golang:1.27-alpine AS builder
ENV GOPROXY=https://goproxy.io,https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o broker-service ./cmd/api

# Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/broker-service .
EXPOSE 8080
CMD ["./broker-service"]