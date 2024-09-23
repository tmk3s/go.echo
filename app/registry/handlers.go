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

func (i *Registry) NewUserHandler() *handler.UserHandler {
	return handler.NewUserHandler(i.NewUserUseCase())
}

func (i *Registry) NewDepartmentHandler() *handler.DepartmentHandler {
	return handler.NewDepartmentHandler(i.NewDepartmentUseCase())
}
