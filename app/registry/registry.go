package registry

import (
	"gorm.io/gorm"
)

// Registry DIコンテナのような依存関係を定義
type Registry struct {
	DbConn *gorm.DB
}

func NewReigistry(dbConn *gorm.DB) Registry {
	return Registry(dbConn)
}

// type Registry struct {}

// func NewReigistry() Registry {
// 	return Registry{}
// }