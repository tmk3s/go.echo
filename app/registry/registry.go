package registry

import (
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// Registry DIコンテナのような依存関係を定義
type Registry struct {
	DbConn      *gorm.DB
	AsynqClient *asynq.Client
}

func NewReigistry(DbConn *gorm.DB, asynqClient *asynq.Client) Registry {
	return Registry{DbConn: DbConn, AsynqClient: asynqClient}
}
