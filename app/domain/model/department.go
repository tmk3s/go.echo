package model

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	CompanyId uint   `json:"company_id" gorm:"index"`
	Name      string `json:"name" gorm:"index"`
	Depth     int    `json:"depth"`
	Ancestors []DepartmentPath `gorm:"foreignKey:Ancestor"`
	Descendants []DepartmentPath `gorm:"foreignKey:Descendant"`
}

func NewDepartment(CompanyId uint, name string) *Department {
	return &Department{
		CompanyId: CompanyId,
		Name:      name,
	}
}
