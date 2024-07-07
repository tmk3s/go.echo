package db

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email	 string `json:"email" gorm:"index`
    Password string `json:"password"`
}

func CreateUser(user *User) {
    db.Create(user)
}

func FindUserById(id uint) User {
    var user User
    db.First(&user, id)
    // db.Where(u).First(&user)
    return user
}

func FindUser(u *User) User {
    var user User
    db.Where(u).First(&user)
    return user
}