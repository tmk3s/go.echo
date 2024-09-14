package registry

import (
	"app/presentation/api/handler"
)

func (i *Registry) NewAuthHandler() *handler.AuthHandler {
	return handler.NewAuthHandler(i.NewAuthUseCase())
}

func (i *Registry) NewTodoHandler() *handler.TodoHandler {
	return handler.NewTodoHandler(i.NewTodoUseCase())
}
