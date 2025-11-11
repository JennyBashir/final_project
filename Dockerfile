FROM golang:1.24 AS builder
LABEL authors="evgeniabashir"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY  . .

RUN go build -o todo-server .

FROM ubuntu:latest

WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/todo-server /app/todo-server
COPY --from=builder /app/web /app/web

RUN mkdir -p /data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

EXPOSE 7540

ENTRYPOINT ["/app/todo-server"]
