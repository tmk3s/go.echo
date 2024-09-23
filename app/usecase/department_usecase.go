package usecase

import (
	"app/domain/model"
	"app/domain/repository"
)

// インターフェースは頭大文字
type DepartmentUseCase interface {
	GetDepartments(companyId uint) (*[]model.Department, error)
	Create(companyId uint, name string, parentId *uint) error
}

type departmentUseCase struct {
	repository.DepartmentRepository
}

func NewDepartmentUseCase(r repository.DepartmentRepository) DepartmentUseCase {
	return &departmentUseCase{r}
}

func (u *departmentUseCase) GetDepartments(companyId uint) (*[]model.Department, error) {
	departments, err := u.DepartmentRepository.GetList(companyId)
	if err != nil {
		return nil, err
	}

	
	return &departments, nil
}

func (u *departmentUseCase) Create(companyId uint, name string, parentId *uint) error {
	department := model.NewDepartment(companyId, name)
	_, err := u.DepartmentRepository.Create(department, parentId)
	if err != nil {
		return err
	}
	return nil
}
