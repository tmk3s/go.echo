package db

import (
	"gorm.io/gorm"
)

type UserAddress struct {
		gorm.Model
    UserId     uint   `json:"user_id" gorm:"praimaly_key"`
    address1     string `json:"address1"`
    address2	 string `json:"address2"`
}
