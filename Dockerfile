FROM golang:1.24.3-alpine3.20 AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

RUN go version

COPY . .

RUN go mod download

RUN go build -o .bin/main cmd/main/main.go

CMD ["make","run"]
