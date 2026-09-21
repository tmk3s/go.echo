package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	CompanyId uint     `json:"company_id" gorm:"index"`
	Email     string   `json:"email" gorm:"size:255;uniqueIndex"`
	Password  string   `json:"-"`
	UserInfo  UserInfo `gorm:"foreignKey:UserId"`
}

func NewUser(email string, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}

// NewCompanyOwner は会社登録時に作成される最初のユーザーを返す。
// CompanyId は会社を作成したあとにリポジトリ側で設定する。
func NewCompanyOwner(email string, password string, lastName string, firstName string) *User {
	return &User{
		Email:    email,
		Password: password,
		UserInfo: UserInfo{
			LastName:  lastName,
			FirstName: firstName,
			Working:   true,
		},
	}
}
