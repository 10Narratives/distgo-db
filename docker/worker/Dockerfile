# Используем официальный Go образ как базовый
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

FROM alpine:latest

COPY --from=builder /app/worker /worker

CMD ["/worker"]