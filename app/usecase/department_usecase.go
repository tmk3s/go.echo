package usecase

import (
	"app/domain/model"
	"app/domain/repository"
)

// インターフェースは頭大文字
type DepartmentUseCase interface {
	GetDepartments(companyId uint) (*[]model.Department, error)
	Create(companyId uint, name string, parentId *uint) error
	Update(id uint, name string) error
	Delete(id uint) error
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

func (u *departmentUseCase) Update(id uint, name string) error {
	department, err := u.DepartmentRepository.GetById(id)
	if err != nil {
		return err
	}
	department.Name = name
	if _, err := u.DepartmentRepository.Update(department); err != nil {
		return err
	}
	return nil
}

func (u *departmentUseCase) Delete(id uint) error {
	department, err := u.DepartmentRepository.GetById(id)
	if err != nil {
		return err
	}
	if err := u.DepartmentRepository.Delete(department); err != nil {
		return err
	}
	return nil
}
