# go.echo
go

https://echo.labstack.com/docs/quick-start

Dockerfile
```
# https://hub.docker.com/_/golang
FROM golang:1.22

ENV LANG C.UTF-8
ENV APP_ROOT /app
WORKDIR /usr/src/app
```

docker-compose.yml
```
version: '3' # composeファイルのバージョン
services:
  app: # サービス名
    build: # ビルドに使うDockerファイルのパス
      context: .
      dockerfile: ./Dockerfile
    volumes: # マウントディレクトリ
      - ./app:/usr/src/app
    tty: true # コンテナの永続化
    # env_file: # .envファイル
    #   - ./build/.go_env
    environment:
      - TZ=Asia/Tokyo
```

docker-cmpose up -d
docker exec -it app /bin/bash

root@d30898d171ec:/usr/src/app# go mod init app
go: creating new go.mod: module app

root@d30898d171ec:/usr/src/app# go version
go version go1.22.4 linux/amd64
root@d30898d171ec:/usr/src/app# go get github.com/labstack/echo/v4
⇩
go.sumが作成される


Create server.go

package main

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

Start server

$ go run server.go

⇩

docker ps
CONTAINER ID   IMAGE        COMMAND   CREATED         STATUS         PORTS                    NAMES
183eb6fcc4c7   goecho-app   "bash"    4 seconds ago   Up 4 seconds   0.0.0.0:1323->1323/tcp   goecho-app-1

port=>OK