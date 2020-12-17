FROM golang:1.15.2-alpine as build

WORKDIR /go/app

COPY src .

RUN apk add --no-cache vim \
    && go build -o app

# 本番環境用マルチステージビルド
FROM alpine

WORKDIR /app

COPY --from=build /go/app/app .
COPY --from=build /go/app/production.env ./.env

CMD ["./app"]
