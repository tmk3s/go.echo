package usecase

import (
	"app/domain/model"
	"app/domain/repository"
)

type AuthUseCase interface {
	GetUser(email string, password string) (*model.User, error)
}

type authUseCase struct {
	repository.UserRepository
}

func NewAuthUseCase(r repository.UserRepository) AuthUseCase {
	return &authUseCase{r}
}

func (u *authUseCase) GetUser(email string, password string) (*model.User, error) {
	return u.UserRepository.GetByEmailAndPass(email, password)
}
