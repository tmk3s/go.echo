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

func (r *companyRepository) Register(company *model.Company, user *model.User) error {
	// 会社だけ作られてユーザーが作られない中途半端な状態を残さないよう1トランザクションにまとめる
	return r.Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(company).Error; err != nil {
			return err
		}
		hashed, err := hashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
		user.CompanyId = company.ID
		// UserInfo は User の関連として同時に作成される
		return tx.Create(user).Error
	})
}
