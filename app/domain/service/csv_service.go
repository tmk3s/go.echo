package service

import (
	"app/domain/model"
	"mime/multipart"
)

type TenureCSVRow struct {
	JoinedOn        string
	ResignationOn   string
	ResignationType string
	Status          string
}

type EmployeeCSVRow struct {
	StaffCode       string
	LastName        string
	FirstName       string
	LastNameKana    string
	FirstNameKana   string
	Email           string
	PostCode        string
	PrefectureName  string
	City            string
	AddressLine1    string
	AddressLine2    string
	Tel             string
	Tenures         []TenureCSVRow
	DepartmentNames []string
}

type CsvService interface {
	ParseDepartmentNames(file multipart.File) ([]string, error)
	GenerateEmployeeCSV(employees []model.Employee) ([]byte, error)
	ParseEmployeeRows(file multipart.File) ([]EmployeeCSVRow, error)
}
