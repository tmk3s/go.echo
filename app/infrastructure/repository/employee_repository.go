package repository

import (
	"errors"

	"app/domain/model"
	"app/domain/repository"

	"gorm.io/gorm"
)

type employeeRepository struct {
	Conn *gorm.DB
}

func NewEmployeeRepository(Conn *gorm.DB) repository.EmployeeRepository {
	return &employeeRepository{Conn}
}

func (r *employeeRepository) GetList(companyId uint) ([]model.Employee, error) {
	var employees []model.Employee
	err := r.Conn.Where(model.Employee{CompanyId: companyId}).Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepository) GetById(companyId uint, id uint) (*model.Employee, error) {
	employee := &model.Employee{}
	err := r.Conn.Where(model.Employee{CompanyId: companyId}).First(employee, id).Error
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (r *employeeRepository) GetAddress(companyId uint, employeeId uint) (*model.EmployeeAddress, error) {
	address := &model.EmployeeAddress{}
	err := r.Conn.
		Select("employee_addresses.*, prefectures.name as prefecture_name").
		Joins("LEFT JOIN prefectures ON employee_addresses.prefecture_id = prefectures.id AND prefectures.deleted_at IS NULL").
		Where("employee_addresses.company_id = ? AND employee_addresses.employee_id = ?", companyId, employeeId).
		Where("employee_addresses.deleted_at IS NULL").
		First(address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return address, nil
}

func (r *employeeRepository) GetTenures(companyId uint, employeeId uint) ([]model.EmployeeTenures, error) {
	var tenures []model.EmployeeTenures
	err := r.Conn.
		Where("company_id = ? AND employee_id = ?", companyId, employeeId).
		Order("joined_on DESC").
		Find(&tenures).Error
	if err != nil {
		return nil, err
	}
	return tenures, nil
}

func (r *employeeRepository) GetDepartments(companyId uint, employeeId uint) ([]model.Department, error) {
	var departments []model.Department
	err := r.Conn.
		Joins("INNER JOIN employee_departments ON departments.id = employee_departments.department_id").
		Where("employee_departments.company_id = ? AND employee_departments.employee_id = ?", companyId, employeeId).
		Where("employee_departments.deleted_at IS NULL").
		Where("departments.deleted_at IS NULL").
		Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *employeeRepository) Create(employee *model.Employee) (*model.Employee, error) {
	if err := r.Conn.Create(employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

func (r *employeeRepository) UpdateEmployee(employee *model.Employee) error {
	return r.Conn.Save(employee).Error
}

func (r *employeeRepository) UpsertAddress(address *model.EmployeeAddress) error {
	var existing model.EmployeeAddress
	err := r.Conn.
		Where("company_id = ? AND employee_id = ?", address.CompanyId, address.EmployeeId).
		Where("deleted_at IS NULL").
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.Conn.Create(address).Error
	}
	if err != nil {
		return err
	}

	address.ID = existing.ID
	address.CreatedAt = existing.CreatedAt
	return r.Conn.Save(address).Error
}

func (r *employeeRepository) UpdateTenure(tenure *model.EmployeeTenures) error {
	return r.Conn.Save(tenure).Error
}

func (r *employeeRepository) UpdateDepartments(companyId uint, employeeId uint, departmentIds []uint) error {
	return r.Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("company_id = ? AND employee_id = ?", companyId, employeeId).
			Delete(&model.EmployeeDepartments{}).Error; err != nil {
			return err
		}
		if len(departmentIds) == 0 {
			return nil
		}
		records := make([]model.EmployeeDepartments, len(departmentIds))
		for i, deptId := range departmentIds {
			records[i] = model.EmployeeDepartments{
				CompanyId:    companyId,
				EmployeeId:   employeeId,
				DepartmentId: deptId,
			}
		}
		return tx.Create(&records).Error
	})
}
