package model

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	CompanyId      uint   `json:"company_id"`
	JobType        string `json:"job_type"`
	Status         string `json:"status"`
	TotalCount     int    `json:"total_count"`
	ProcessedCount int    `json:"processed_count"`
	ErrorMessage   string `json:"error_message"`
	FileData       []byte `json:"-"`
}
