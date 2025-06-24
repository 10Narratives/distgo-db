FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o master ./cmd/master

FROM alpine:latest

COPY --from=builder /app/master /master

EXPOSE 50051

CMD ["/master"]