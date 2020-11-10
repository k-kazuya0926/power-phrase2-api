FROM golang:1.15.4-alpine as build

WORKDIR /go/src/github.com/k-kazuya0926/power-phrase2-api

COPY src .

RUN apk update \
    && apk add --no-cache vim
