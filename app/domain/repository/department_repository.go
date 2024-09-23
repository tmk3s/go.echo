package repository

import "app/domain/model"

// 実装はapp/infrastructure/repository/todo_repository.goにて行う
type DepartmentRepository interface {
	GetList(companyId uint) ([]model.Department, error)
	Create(department *model.Department, parentId *uint) (*model.Department, error)
}
