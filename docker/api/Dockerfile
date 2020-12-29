FROM golang:1.15.2-alpine as build

WORKDIR /go/app

COPY src .

RUN set -eux && \
    apk update && \
    apk add --no-cache less && \
    go get -u github.com/cosmtrek/air && \
    go build -o /go/bin/air github.com/cosmtrek/air && \
    go get -u github.com/go-delve/delve/cmd/dlv && \
    go build -o /go/bin/dlv github.com/go-delve/delve/cmd/dlv && \
    go build -o app

# 本番環境用マルチステージビルド
FROM alpine

WORKDIR /app

COPY --from=build /go/app/app .

CMD ["./app"]
