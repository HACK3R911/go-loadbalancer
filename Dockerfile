FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/loadbalancer ./cmd/loadbalancer/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/loadbalancer ./loadbalancer
COPY config.yaml /app/config.yaml
EXPOSE 8080
CMD ["./loadbalancer", "-config=/app/config.yaml"]  