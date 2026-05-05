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

func (i *Registry) NewUserUseCase() usecase.UserUseCase {
	return usecase.NewUserUseCase(
		i.NewUserRepository(),
	)
}

func (i *Registry) NewDepartmentUseCase() usecase.DepartmentUseCase {
	return usecase.NewDepartmentUseCase(
		i.NewDepartmentRepository(),
		i.NewCsvService(),
	)
}

func (i *Registry) NewEmployeeUseCase() usecase.EmployeeUseCase {
	return usecase.NewEmployeeUseCase(i.NewEmployeeRepository())
}

func (i *Registry) NewPrefectureUseCase() usecase.PrefectureUseCase {
	return usecase.NewPrefectureUseCase(i.NewPrefectureRepository())
}
