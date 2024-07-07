package db

import (
	"gorm.io/gorm"
)

type Tag struct {
		gorm.Model
    Name     string `json:"name" gorm:"index`
}
