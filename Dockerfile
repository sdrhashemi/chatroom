FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o chatroom ./cmd/server/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/chatroom .

COPY ./config.yaml .

EXPOSE 8080




