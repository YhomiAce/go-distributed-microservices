# Build
FROM golang:1.26-alpine AS builder
# ENV GOPROXY=https://goproxy.io,https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/api

# Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service .
EXPOSE 8081
CMD ["./auth-service"]