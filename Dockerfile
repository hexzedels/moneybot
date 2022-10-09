# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY .env ./
COPY go.mod ./
COPY go.sum ./

RUN apk update && \
    apk add build-base

RUN go mod download

COPY *.go ./

RUN go build -o ./money-bot

EXPOSE 8080

CMD [ "./money-bot" ]
