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

https://echo.labstack.com/docs/middleware/jwt
https://echo.labstack.com/docs/middleware/logger#configuration
https://qiita.com/atsutama/items/68773112208c7fc05069
https://qiita.com/x-color/items/24ff2491751f55e866cf
https://echo.labstack.com/docs/binding
https://gorm.io/ja_JP/docs/delete.html
https://apidog.com/jp/blog/axios-put-request/
https://github.com/golang-jwt/jwt/blob/main/example_test.go

docker でログ見れない問題
エラーチェック毎回しないといけない
ポインタとか
jwt周りの設定



### コマンド

#### サーバー起動
```bash
docker compose up -d
```

#### DB初期化 & Seed投入
DBをdrop/createしてマイグレーション・Seedを一括実行する。
```bash
docker compose exec backend go run ./cmd/refresh/main.go
```

#### Seedのみ投入
```bash
docker compose exec backend go run ./cmd/seed/main.go
```

#### テスト実行
全パッケージのテストを実行する。
```bash
docker compose exec backend go test ./...
```

パッケージ単位で実行する場合：
```bash
# usecase
docker compose exec backend go test ./usecase/...

# repository（MySQL の develop_test DB を使用）
docker compose exec backend go test ./infrastructure/repository/...

# CSV service
docker compose exec backend go test ./infrastructure/service/...
```

特定のテストのみ実行する場合：
```bash
docker compose exec backend go test ./usecase/... -run TestBulkCreateFromCSV
```

テスト概要：

| 層 | 種別 | 備考 |
|---|---|---|
| `usecase/` | ユニットテスト（モック） | DB 接続不要 |
| `infrastructure/repository/` | インテグレーションテスト | MySQL（develop_test DB）を使用 |
| `infrastructure/service/` | ユニットテスト | DB 接続不要 |

### ディレクトリ構造
- domain
  - model
  - repository
  - service
- errors
- frontend
- infrastructure
  - dto
  - repository
  - service
- presentation
  - api
    - handler
    - router
- usecase

#### domain/model
ドメインを表現する(railsのmodelと同じ)

#### domain/repository
interface群
infrastructure>repositoryで実装する

#### domain/service
ドメインのビジネスロジックで単体のモデルでは表現できないものはここに定義する
interface群
infrastructure>serviceで実装する

### infrastructure>repository
repositoryの実装を行う

### infrastructure>service
serviceの実装を行う

### presentation>api>handler
APIの定義

### presentation>api>router
APIのパス定義

### usecase
presentation>api>handlerから呼ばれる

---

## Claude Code コードレビュー

`.claude/` に、複数の専門サブエージェントへ並列でレビューを投げて結果を 1 本のレポートに統合する仕組みを用意している。

```
.claude/
├── skills/review/SKILL.md        # /review の本体（対象の判定・並列起動・結果の統合）
└── agents/
    ├── go-backend-reviewer.md    # Go / Echo / GORM / クリーンアーキテクチャ
    ├── frontend-reviewer.md      # Next.js(App Router) / TypeScript / Tailwind
    └── security-reviewer.md      # 認証認可・入力検証・機密情報（前2者を横断）
```

### 使い方

Claude Code のセッションで以下を実行する。

| コマンド | モード | 対象 |
|---|---|---|
| `/review` | 差分 | 未コミット差分（無ければ main との差分） |
| `/review staged` | 差分 | ステージ済みの変更 |
| `/review #42` | 差分 | PR #42（`gh pr diff` を使用） |
| `/review feature/xxx` | 差分 | main と `feature/xxx` の差分 |
| `/review all` | 全体 | `app/` 配下のファイル全体 |
| `/review app/usecase` | 全体 | 指定パス配下のファイル全体 |

- **差分モード**: 変更行とその影響範囲だけを見る。コミット前・PR 作成前の確認向け。
- **全体モード**: 指定ファイルを最初から最後まで読む。行単位の指摘に加えて、レイヤー責務のずれ・重複処理・設計の一貫性といったファイル横断の観点も報告される。

### 動作

1. 引数からモードと対象ファイルを決定する（対象が空なら何も起動せず終了）
2. 変更パスに応じて起動するエージェントを選ぶ
   - `app/**/*.go` があれば `go-backend-reviewer`
   - `app/frontend/**` があれば `frontend-reviewer`
   - `security-reviewer` は常に起動
3. 選んだエージェントを **並列** で起動する
4. 各エージェントの指摘を統合し、重複をまとめ、矛盾があれば実際のコードを確認して裁定、誤検出を除外
5. severity（high → low）順に並べた 1 本のレポートを出力する

レビューは指摘までで、**コードの自動修正は行わない**。修正が必要な場合はレポートを見てから別途依頼する。

### 注意

- `/review all` は `app/` 全体が対象になりファイル数が多い。40 ファイルを超える場合は実行前に確認が入るが、まずは `/review app/usecase` のようにパスを絞るのが速い。
- `.claude/agents/*.md` を新規作成・変更した場合、エージェント定義はセッション起動時に読み込まれるため、**Claude Code のセッションを再起動**してから実行する。

### カスタマイズ

- レビュー観点を足したい → 該当する `.claude/agents/*.md` の「重点チェック項目」に追記する
- 新しい観点のエージェントを追加したい → `.claude/agents/` に同じ形式で md を追加し、`SKILL.md` のステップ 2 の振り分け表に行を足す
- レポートの体裁を変えたい → `SKILL.md` のステップ 5 のテンプレートを編集する
