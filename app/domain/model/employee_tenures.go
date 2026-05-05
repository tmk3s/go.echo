package model

import (
	"time"

	"gorm.io/gorm"
)

type EmployeeTenures struct {
	gorm.Model
	CompanyId      uint       `json:"company_id" gorm:"index"`
	EmployeeId     uint       `json:"employee_id" gorm:"index"`
	JoinedOn       time.Time  `json:"joined_on"`
	ResignationOn  *time.Time `json:"resignation_on"`
	ResignationType *string   `json:"resignation_type"`
	Status         *string    `json:"status"`
}
