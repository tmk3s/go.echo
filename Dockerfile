# https://hub.docker.com/_/golang
FROM golang:1.22

ENV LANG C.UTF-8
ENV APP_ROOT /app
WORKDIR /usr/src/app

RUN apt-get update -qq && apt-get install -y vim && apt-get install -y nodejs npm

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY app/go.mod app/go.sum ./
RUN go mod download && go mod verify
RUN go install github.com/air-verse/air@latest

# COPY . .
# RUN go build -v -o /usr/local/bin/app ./...

# CMD ["app"]