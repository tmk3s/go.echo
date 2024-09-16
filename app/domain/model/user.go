package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	CompanyId uint   `json:"company_id" gorm:"index"`
	Email     string `json:"email" gorm:"index`
	Password  string `json:"password"`
	UserInfo  UserInfo `gorm:"foreignKey:UserId"`
}

func NewUser(email string, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}
