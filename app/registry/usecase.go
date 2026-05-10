package registry

import (
	"app/infrastructure/worker"
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
	return usecase.NewEmployeeUseCase(
		i.NewEmployeeRepository(),
		i.NewDepartmentRepository(),
		i.NewPrefectureRepository(),
		i.NewCsvService(),
		i.NewJobRepository(),
		worker.NewAsynqEnqueuer(i.AsynqClient),
	)
}

func (i *Registry) NewJobUseCase() usecase.JobUseCase {
	return usecase.NewJobUseCase(i.NewJobRepository())
}

func (i *Registry) NewPrefectureUseCase() usecase.PrefectureUseCase {
	return usecase.NewPrefectureUseCase(i.NewPrefectureRepository())
}

func (i *Registry) NewCompanyUseCase() usecase.CompanyUseCase {
	return usecase.NewCompanyUseCase(i.NewCompanyRepository())
}
