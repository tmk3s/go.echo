package service

import "mime/multipart"

type CsvService interface {
	ParseDepartmentNames(file multipart.File) ([]string, error)
}
