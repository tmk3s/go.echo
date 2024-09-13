package registry

import (
	"app/usecase"
)

func (i *Registry) NewAuthUseCase() usecase.AuthUseCase {
	return usecase.NewAuthUseCase(i.NewUserRepository())
}

func (i *Registry) NewTodoUseCase() usecase.TodoUseCase {
	return usecase.NewTodoUseCase(
		i.NewTodoRepository(),
		i.NewUserRepository(),
	)
}