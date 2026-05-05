package model

import (
	"gorm.io/gorm"
)

type EmployeeAddress struct {
	gorm.Model
	CompanyId      uint    `json:"company_id" gorm:"index"`
	EmployeeId     uint    `json:"employee_id" gorm:"index"`
	PostCode       *string `json:"post_code"`
	PrefectureId   *uint   `json:"prefecture_id"`
	// gorm:"-" によりDBカラムは作られず、JOINクエリで取得した値を手動でマッピングするためだけに使用
	PrefectureName *string `json:"prefecture_name" gorm:"-"`
	City           *string `json:"city"`
	AddressLine1   *string `json:"address_line1"`
	AddressLine2   *string `json:"address_line2"`
	Tel            *string `json:"tel"`
}
