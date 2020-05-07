FROM golang:1.13-alpine

RUN apk update && apk add git

WORKDIR /go/src/github.com/dcrichards/go-to-openapi

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
