package usecase

import (
	"mime/multipart"

	"app/domain/model"
	"app/domain/repository"
	"app/domain/service"
)

type DepartmentUseCase interface {
	GetDepartments(companyId uint) (*[]model.Department, error)
	Create(companyId uint, name string, parentId *uint) error
	Update(id uint, name string) error
	Delete(id uint) error
	Download() error
	Upload(companyId uint, file multipart.File, fileHeader *multipart.FileHeader) error
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

func (u *departmentUseCase) Upload(companyId uint, file multipart.File, fileHeader *multipart.FileHeader) error {
	names, err := u.csvService.ParseDepartmentNames(file)
	if err != nil {
		return err
	}

	existing, err := u.repo.GetList(companyId)
	if err != nil {
		return err
	}

	existingNames := make(map[string]bool)
	for _, d := range existing {
		existingNames[d.Name] = true
	}

	for _, name := range names {
		if !existingNames[name] {
			department := model.NewDepartment(companyId, name)
			if _, err := u.repo.Create(department, nil); err != nil {
				return err
			}
		}
	}
	return nil
}
