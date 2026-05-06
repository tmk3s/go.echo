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
	return employees, err
}

func (r *employeeRepository) GetDetail(companyId uint, id uint) (*model.Employee, error) {
	var emp model.Employee
	err := r.Conn.
		Preload("Address").
		Preload("Tenures", func(db *gorm.DB) *gorm.DB {
			return db.Order("joined_on DESC")
		}).
		Where("company_id = ?", companyId).
		First(&emp, id).Error
	if err != nil {
		return nil, err
	}
	r.resolvePrefectureName(emp.Address)

	depts, err := r.loadDepartments(companyId, emp.ID)
	if err != nil {
		return nil, err
	}
	emp.Departments = depts
	return &emp, nil
}

func (r *employeeRepository) GetListForExport(companyId uint) ([]model.Employee, error) {
	var employees []model.Employee
	err := r.Conn.
		Preload("Address").
		Preload("Tenures", func(db *gorm.DB) *gorm.DB {
			return db.Order("joined_on DESC")
		}).
		Where("company_id = ?", companyId).
		Find(&employees).Error
	if err != nil {
		return nil, err
	}

	// 都道府県名をまとめて解決（N+1防止）
	var prefects []model.Prefecture
	if err := r.Conn.Find(&prefects).Error; err != nil {
		return nil, err
	}
	prefMap := make(map[uint]string, len(prefects))
	for _, p := range prefects {
		prefMap[p.ID] = p.Name
	}

	// 部署をまとめて解決（N+1防止）
	deptMap, err := r.loadAllDepartments(companyId)
	if err != nil {
		return nil, err
	}

	for i := range employees {
		if employees[i].Address != nil && employees[i].Address.PrefectureId != nil {
			name := prefMap[*employees[i].Address.PrefectureId]
			employees[i].Address.PrefectureName = &name
		}
		employees[i].Departments = deptMap[employees[i].ID]
	}
	return employees, nil
}

func (r *employeeRepository) loadDepartments(companyId uint, employeeId uint) ([]model.Department, error) {
	var departments []model.Department
	err := r.Conn.
		Joins("INNER JOIN employee_departments ON departments.id = employee_departments.department_id").
		Where("employee_departments.company_id = ? AND employee_departments.employee_id = ?", companyId, employeeId).
		Where("employee_departments.deleted_at IS NULL").
		Where("departments.deleted_at IS NULL").
		Find(&departments).Error
	return departments, err
}

func (r *employeeRepository) loadAllDepartments(companyId uint) (map[uint][]model.Department, error) {
	type deptRow struct {
		ID         uint
		Name       string
		Depth      int
		OrderNo    int
		EmployeeId uint
	}
	var rows []deptRow
	err := r.Conn.
		Table("departments").
		Select("departments.id, departments.name, departments.depth, departments.order_no, employee_departments.employee_id").
		Joins("INNER JOIN employee_departments ON departments.id = employee_departments.department_id").
		Where("employee_departments.company_id = ? AND employee_departments.deleted_at IS NULL AND departments.deleted_at IS NULL", companyId).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[uint][]model.Department)
	for _, row := range rows {
		dept := model.Department{Name: row.Name, Depth: row.Depth, OrderNo: row.OrderNo}
		dept.ID = row.ID
		result[row.EmployeeId] = append(result[row.EmployeeId], dept)
	}
	return result, nil
}

func (r *employeeRepository) resolvePrefectureName(addr *model.EmployeeAddress) {
	if addr == nil || addr.PrefectureId == nil {
		return
	}
	var pref model.Prefecture
	if err := r.Conn.First(&pref, *addr.PrefectureId).Error; err == nil {
		addr.PrefectureName = &pref.Name
	}
}

func (r *employeeRepository) Create(employee *model.Employee) (*model.Employee, error) {
	if err := r.Conn.Create(employee).Error; err != nil {
		return nil, err
	}
	return employee, nil
}

func (r *employeeRepository) BulkCreate(employees []model.Employee) error {
	return r.Conn.CreateInBatches(employees, 100).Error
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

func (r *employeeRepository) ReplaceTenures(companyId uint, employeeId uint, tenures []model.EmployeeTenures) error {
	return r.Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("company_id = ? AND employee_id = ?", companyId, employeeId).
			Delete(&model.EmployeeTenures{}).Error; err != nil {
			return err
		}
		if len(tenures) == 0 {
			return nil
		}
		return tx.Create(&tenures).Error
	})
}
