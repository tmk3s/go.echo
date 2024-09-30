package repository

import "app/domain/model"

type DepartmentRepository interface {
	GetById(id uint) (*model.Department, error)
	GetList(companyId uint) ([]model.Department, error)
	Create(department *model.Department, parentId *uint) (*model.Department, error)
	Update(department *model.Department) (*model.Department, error)
	Delete(department *model.Department) error
}
