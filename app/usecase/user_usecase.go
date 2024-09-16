package usecase

import (
	"time"

	"app/domain/model"
	"app/domain/repository"
)

type UserParams struct {
	last_name  string
	first_name string
	email      string
	gender     int
	birth_day  time.Time
	Working    bool
}

// インターフェースは頭大文字
type UserUseCase interface {
	Get(userId uint) (*model.User, error)
	Update(userId uint, params UserParams) error
}

type userUseCase struct {
	repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) UserUseCase {
	return &userUseCase{r}
}

func (u *userUseCase) Get(userId uint) (*model.User, error) {
	user, err := u.UserRepository.GetById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) Update(userId uint, params UserParams) error {
	_, err := u.UserRepository.GetById(userId)
	if err != nil {
		return err
	}
	return nil
}
