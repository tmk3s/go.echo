package usecase

import (
	"app/domain/model"
	"app/domain/repository"
)

type CompanyUseCase interface {
	GetCompany(id uint) (*model.Company, error)
	UpdateCompany(id uint, name string) error
}

type companyUseCase struct {
	repo repository.CompanyRepository
}

func NewCompanyUseCase(r repository.CompanyRepository) CompanyUseCase {
	return &companyUseCase{repo: r}
}

func (u *companyUseCase) GetCompany(id uint) (*model.Company, error) {
	return u.repo.GetById(id)
}

func (u *companyUseCase) UpdateCompany(id uint, name string) error {
	company, err := u.repo.GetById(id)
	if err != nil {
		return err
	}
	company.Name = name
	return u.repo.Update(company)
}
