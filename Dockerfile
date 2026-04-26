# STAGE 1: Compilación
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gopherguard ./cmd/server/main.go

# STAGE 2: Imagen Final Pro
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /gopherguard /gopherguard
# Copiamos las propiedades por ambiente
COPY vars/ /vars/
EXPOSE 8080
ENTRYPOINT ["/gopherguard"]