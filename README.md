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

ホットリロード
https://qiita.com/frkawa/items/1dd12e19f10de034e0f5
https://github.com/air-verse/air/issues/605#issuecomment-2146474109

コマンド： go install github.com/air-verse/air@latest



2024-06-09 02:20:34 presentation/api/handler/auth_handler.go has changed
2024-06-09 02:20:34 building...
2024-06-09 02:20:34 db/db.go:4:5: no required module provides package github.com/jinzhu/gorm; to add it:

2024-06-09 02:20:34     go get github.com/jinzhu/gorm

2024-06-09 02:20:34 db/db.go:5:5: no required module provides package github.com/jinzhu/gorm/dialects/sqlite; to add it:

2024-06-09 02:20:34     go get github.com/jinzhu/gorm/dialects/sqlite

2024-06-09 02:20:34 presentation/api/handler/auth_handler.go:7:5: no required module provides package github.com/dgrijalva/jwt-go; to add it:

2024-06-09 02:20:34     go get github.com/dgrijalva/jwt-go

2024-06-09 02:20:34 presentation/api/handler/auth_handler.go:8:5: no required module provides package github.com/labstack/echo; to add it:

2024-06-09 02:20:34     go get github.com/labstack/echo

2024-06-09 02:20:34 presentation/api/handler/auth_handler.go:9:5: no required module provides package github.com/labstack/echo/middleware; to add it:

2024-06-09 02:20:34     go get github.com/labstack/echo/middleware
⇩
https://sumito.jp/2021/04/23/go1-16-build-error-github-com-missing-package/
go mod tidy