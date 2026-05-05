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

func (i *Registry) NewDepartmentRepository() repository.DepartmentRepository {
	return repositoryImpl.NewDepartmentRepository(i.DbConn)
}

func (i *Registry) NewEmployeeRepository() repository.EmployeeRepository {
	return repositoryImpl.NewEmployeeRepository(i.DbConn)
}

func (i *Registry) NewPrefectureRepository() repository.PrefectureRepository {
	return repositoryImpl.NewPrefectureRepository(i.DbConn)
}
