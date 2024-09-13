package repository

import "app/domain/model"

type UserRepository interface {
	GetEmailAndPass(email string, password string) (*model.User, error)
	Create(email string, password string) (*model.User, error)
	Update(u *model.User) (*model.User, error)
	Delete(id uint) error
}