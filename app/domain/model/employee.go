package model

import (
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	CompanyId     uint              `json:"company_id" gorm:"index"`
	LastName      string            `json:"last_name"`
	FirstName     string            `json:"first_name"`
	LastNameKana  string            `json:"last_name_kana"`
	FirstNameKana string            `json:"first_name_kana"`
	Email         string            `json:"email" gorm:"index"`
	StaffCode     string            `json:"staff_code" gorm:"index"`
	Address       *EmployeeAddress  `json:"address,omitempty" gorm:"foreignKey:EmployeeId"`
	Tenures       []EmployeeTenures `json:"tenures,omitempty" gorm:"foreignKey:EmployeeId"`
	Departments   []Department      `json:"departments,omitempty" gorm:"-"`
}
