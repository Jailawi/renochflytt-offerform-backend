FROM golang:1.24 AS builder

WORKDIR /app

COPY cmd /app/cmd
COPY pkg /app/pkg
COPY database /app/database
COPY templates /app/templates
COPY utils /app/utils
COPY go.mod go.sum /app/

RUN go mod download
RUN go build -o app ./cmd/main.go

ENTRYPOINT [ "./app" ]
