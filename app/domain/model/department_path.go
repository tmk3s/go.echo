package model

import (
	"gorm.io/gorm"
)

type DepartmentPath struct {
	gorm.Model
	AncestorId   uint `json:"ancestor_id" gorm:"index"`
	DescendantId uint `json:"descendant_id" gorm:"index"`
}
