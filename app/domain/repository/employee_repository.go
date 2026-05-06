package repository

import "app/domain/model"

type EmployeeRepository interface {
	GetList(companyId uint) ([]model.Employee, error)
	GetDetail(companyId uint, id uint) (*model.Employee, error)
	GetListForExport(companyId uint) ([]model.Employee, error)
	Create(employee *model.Employee) (*model.Employee, error)
	BulkCreate(employees []model.Employee) error
	UpdateEmployee(employee *model.Employee) error
	UpsertAddress(address *model.EmployeeAddress) error
	UpdateTenure(tenure *model.EmployeeTenures) error
	UpdateDepartments(companyId uint, employeeId uint, departmentIds []uint) error
	ReplaceTenures(companyId uint, employeeId uint, tenures []model.EmployeeTenures) error
}
