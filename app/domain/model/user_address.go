package model

import (
	"gorm.io/gorm"
)

type UserAddress struct {
	gorm.Model
	UserId   uint   `json:"user_id" gorm:"praimaly_key"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
}
