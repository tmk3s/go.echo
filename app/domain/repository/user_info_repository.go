package repository

import "app/domain/model"

type UserInfoRepository interface {
	Create(u *model.UserInfo) (*model.UserInfo, error)
	Update(u *model.UserInfo) (*model.UserInfo, error)
}
