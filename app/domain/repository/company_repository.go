package repository

import "app/domain/model"

type CompanyRepository interface {
	GetById(id uint) (*model.Company, error)
	Update(company *model.Company) error
	// Register は会社と、その最初のユーザーを1つのトランザクションで作成する。
	// user.CompanyId は作成された会社のIDで上書きされる。
	Register(company *model.Company, user *model.User) error
}
