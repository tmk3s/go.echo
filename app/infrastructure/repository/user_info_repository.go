package repository

import "app/domain/model"

type UserRepository interface {
	Upsert(u *model.User) (*model.User, error)
}
