package model

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	CompanyId   uint             `json:"company_id" gorm:"index"`
	Name        string           `json:"name" gorm:"index"`
	Depth       int              `json:"depth"`
	OrderNo     int              `json:"order_no"`
	Ancestors   []DepartmentPath `gorm:"foreignKey:AncestorId"`
	Descendants []DepartmentPath `gorm:"foreignKey:DescendantId"`
}

func NewDepartment(CompanyId uint, name string) *Department {
	return &Department{
		CompanyId: CompanyId,
		Name:      name,
	}
}
