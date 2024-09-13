package registry

import (
	"gorm.io/gorm"
)

// Registry DIコンテナのような依存関係を定義
type Registry struct {
	DbConn *gorm.DB
}

func NewReigistry(DbConn *gorm.DB) Registry {
	return Registry{DbConn}
}