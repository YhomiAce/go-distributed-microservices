# Build
FROM golang:1.27-alpine AS builder
ENV GOPROXY=https://goproxy.io,https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mailer-service ./cmd/api

# Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mailer-service .
EXPOSE 8083
CMD ["./mailer-service"]