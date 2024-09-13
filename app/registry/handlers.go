package registry

import (
	"app/presentation/api/handler"
)

func (i *Registry) NewAuthHandler() *handler.AuthHandler {
	return handler.NewAuthHandler(i.NewAuthUseCase)
}