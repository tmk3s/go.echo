package db

import (
  "time"

	"gorm.io/gorm"
)

type UserInfo struct {
		gorm.Model
    UserId     uint   `json:"user_id" gorm:"praimaly_key"`
    LastName    string `json:"last_name" gorm:"index`
    FirstName   string `json:"first_name"`
    Gender      string `json:"gender"`
    BirthDay    time.Time `json:"birthday"`
    Working     bool      `json:"Working"`
    Image       []byte    `json:"image"`
}
