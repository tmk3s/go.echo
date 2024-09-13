package registry

import (
	"app/domain/repository"
	repositoryImpl "app/infrastructure/repository"
)

func (i *Registry) NewUserRepository() repository.UserRepository {
	return repositoryImpl.NewUserRepository(i.DbConn)
}

func (i *Registry) NewTodoRepository() repository.TodoRepository {
	return repositoryImpl.NewTodoRepository(i.DbConn)
}