# https://hub.docker.com/_/golang
# alpineベースに変更（イメージを軽量化し、脆弱性を減らすため）
FROM golang:1.26-alpine

ENV LANG C.UTF-8
ENV APP_ROOT /app
WORKDIR /usr/src/app

# apk（Alpine）に変更。build-baseはgorm.io/driver/sqliteがcgoを使う（Cコンパイラが必要）ため追加
RUN apk update && apk add --no-cache vim nodejs npm build-base

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY app/go.mod app/go.sum ./
RUN go mod download && go mod verify
RUN go install github.com/air-verse/air@latest

# COPY . .
# RUN go build -v -o /usr/local/bin/app ./...

# CMD ["app"]