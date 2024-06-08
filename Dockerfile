# https://hub.docker.com/_/golang
FROM golang:1.22

ENV LANG C.UTF-8
ENV APP_ROOT /app
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify && go mod init app && go get github.com/labstack/echo/v4

# COPY . .
# RUN go build -v -o /usr/local/bin/app ./...

# CMD ["app"]