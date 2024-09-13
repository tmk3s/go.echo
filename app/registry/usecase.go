package registry

import (
	"app/usecase"
)

func (i *Registry) NewAuthUseCase() *usecase.AuthUseCase {
	return usecase.AuthUseCase(i.NewUserRepository)
}