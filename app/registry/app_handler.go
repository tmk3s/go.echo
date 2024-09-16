package registry

import (
	"app/presentation/api/handler"
)

// handlerのリスト。新規にhandlerを追加するたびにこちらにも追加
func (i *Registry) NewAppHandler() *handler.AppHandler {
	appHandler := &handler.AppHandler{
		AuthHandler: *i.NewAuthHandler(),
		TodoHandler: *i.NewTodoHandler(),
		UserHandler: *i.NewUserHandler(),
	}
	return appHandler
}
