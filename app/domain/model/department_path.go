package model

import (
	"gorm.io/gorm"
)

type DepartmentPath struct {
	gorm.Model
	Ancestor   uint `json:"ancestor" gorm:"index"`
	Descendant uint `json:"descendant" gorm:"index"`
}
