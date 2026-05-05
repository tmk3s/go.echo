package repository

import "app/domain/model"

type EmployeeRepository interface {
	GetList(companyId uint) ([]model.Employee, error)
	GetById(companyId uint, id uint) (*model.Employee, error)
	GetAddress(companyId uint, employeeId uint) (*model.EmployeeAddress, error)
	GetTenures(companyId uint, employeeId uint) ([]model.EmployeeTenures, error)
	GetDepartments(companyId uint, employeeId uint) ([]model.Department, error)
	Create(employee *model.Employee) (*model.Employee, error)
	BulkCreate(employees []model.Employee) error
	UpdateEmployee(employee *model.Employee) error
	UpsertAddress(address *model.EmployeeAddress) error
	UpdateTenure(tenure *model.EmployeeTenures) error
	UpdateDepartments(companyId uint, employeeId uint, departmentIds []uint) error
	// Batch fetch for N+1 free CSV export
	GetAllAddresses(companyId uint) (map[uint]*model.EmployeeAddress, error)
	GetAllTenures(companyId uint) (map[uint][]model.EmployeeTenures, error)
	GetAllDepartments(companyId uint) (map[uint][]model.Department, error)
	// For CSV import
	ReplaceTenures(companyId uint, employeeId uint, tenures []model.EmployeeTenures) error
}
