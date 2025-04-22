#запускать из корня
FROM golang:1.23.3-alpine AS builder

COPY . /github.com/Ippolid/chat-server/source/
WORKDIR /github.com/Ippolid/chat-server/source/

RUN go mod download
RUN go build -o ./bin/server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Ippolid/chat-server/source/bin/server .
COPY ./deploy/.env .env

CMD ["./server"]