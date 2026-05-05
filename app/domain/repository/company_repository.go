package repository

import "app/domain/model"

type CompanyRepository interface {
	GetById(id uint) (*model.Company, error)
	Update(company *model.Company) error
}
