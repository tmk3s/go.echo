package model

import (
	"gorm.io/gorm"
)

type EmployeeDepartments struct {
	gorm.Model
	CompanyId    uint `json:"company_id" gorm:"index"`
	EmployeeId   uint `json:"employee_id" gorm:"index"`
	DepartmentId uint `json:"department_id" gorm:"index"`
}
