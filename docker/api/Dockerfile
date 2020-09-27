FROM golang:1.13.5-alpine as build

WORKDIR /go/src/github.com/k-kazuya0926/power-phrase2-api

COPY src .

RUN apk update \
    && apk add --no-cache vim \
    && go get gopkg.in/urfave/cli.v2@master \
    && go get github.com/oxequa/realize
