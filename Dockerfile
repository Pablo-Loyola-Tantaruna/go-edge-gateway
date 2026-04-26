FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gateway cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gateway .
COPY --from=builder /app/config.yaml .
CMD ["./gateway"]