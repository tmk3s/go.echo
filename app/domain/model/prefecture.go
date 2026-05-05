package model

import (
	"gorm.io/gorm"
)

type Prefecture struct {
	gorm.Model
	Name string `json:"name" gorm:"index"`
}
