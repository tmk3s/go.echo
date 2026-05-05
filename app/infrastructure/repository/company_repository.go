package repository

import (
	"app/domain/model"
	"app/domain/repository"

	"gorm.io/gorm"
)

type companyRepository struct {
	Conn *gorm.DB
}

func NewCompanyRepository(Conn *gorm.DB) repository.CompanyRepository {
	return &companyRepository{Conn}
}

func (r *companyRepository) GetById(id uint) (*model.Company, error) {
	company := &model.Company{}
	if err := r.Conn.First(company, id).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (r *companyRepository) Update(company *model.Company) error {
	return r.Conn.Save(company).Error
}
