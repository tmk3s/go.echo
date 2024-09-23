package model

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	CompanyId uint   `json:"company_id" gorm:"index"`
	Name      string `json:"name" gorm:"index"`
}

func NewDepartment(CompanyId uint, name string) *Department {
	return &Department{
		CompanyId: CompanyId,
		Name:      name,
	}
}
