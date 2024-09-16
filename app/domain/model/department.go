package model

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	CompanyId uint   `json:"company_id" gorm:"index"`
	Name      string `json:"name" gorm:"index"`
}
