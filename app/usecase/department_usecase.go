package usecase

import (
	"app/domain/model"
	"app/domain/repository"
	"app/domain/service"
)

// インターフェースは頭大文字
type DepartmentUseCase interface {
	GetDepartments(companyId uint) (*[]model.Department, error)
	Create(companyId uint, name string, parentId *uint) error
	Update(id uint, name string) error
	Delete(id uint) error
	Download() error
	Upload() error
}

type departmentUseCase struct {
	repo       repository.DepartmentRepository
	csvService service.CsvService
}

func NewDepartmentUseCase(r repository.DepartmentRepository, s service.CsvService) DepartmentUseCase {
	return &departmentUseCase{repo: r, csvService: s}
}

func (u *departmentUseCase) GetDepartments(companyId uint) (*[]model.Department, error) {
	departments, err := u.repo.GetList(companyId)
	if err != nil {
		return nil, err
	}

	return &departments, nil
}

func (u *departmentUseCase) Create(companyId uint, name string, parentId *uint) error {
	department := model.NewDepartment(companyId, name)
	_, err := u.repo.Create(department, parentId)
	if err != nil {
		return err
	}
	return nil
}

func (u *departmentUseCase) Update(id uint, name string) error {
	department, err := u.repo.GetById(id)
	if err != nil {
		return err
	}
	department.Name = name
	if _, err := u.repo.Update(department); err != nil {
		return err
	}
	return nil
}

func (u *departmentUseCase) Delete(id uint) error {
	department, err := u.repo.GetById(id)
	if err != nil {
		return err
	}
	if err := u.repo.Delete(department); err != nil {
		return err
	}
	return nil
}

func (u *departmentUseCase) Download() error {
	return nil
}

func (u *departmentUseCase) Upload() error {
	return nil
}
