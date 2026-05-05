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
	// gorm:"-" フィールドは JOIN でもスキャンされないため、prefecture_name は別クエリで取得する
	err := r.Conn.
		Where("company_id = ? AND employee_id = ?", companyId, employeeId).
		Where("deleted_at IS NULL").
		First(address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if address.PrefectureId != nil {
		var prefecture model.Prefecture
		if err := r.Conn.First(&prefecture, *address.PrefectureId).Error; err == nil {
			address.PrefectureName = &prefecture.Name
		}
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

func (r *employeeRepository) GetAllAddresses(companyId uint) (map[uint]*model.EmployeeAddress, error) {
	var addresses []model.EmployeeAddress
	if err := r.Conn.Where("company_id = ? AND deleted_at IS NULL", companyId).Find(&addresses).Error; err != nil {
		return nil, err
	}
	var prefects []model.Prefecture
	if err := r.Conn.Find(&prefects).Error; err != nil {
		return nil, err
	}
	prefMap := make(map[uint]string, len(prefects))
	for _, p := range prefects {
		prefMap[p.ID] = p.Name
	}
	result := make(map[uint]*model.EmployeeAddress, len(addresses))
	for i := range addresses {
		if addresses[i].PrefectureId != nil {
			name := prefMap[*addresses[i].PrefectureId]
			addresses[i].PrefectureName = &name
		}
		result[addresses[i].EmployeeId] = &addresses[i]
	}
	return result, nil
}

func (r *employeeRepository) GetAllTenures(companyId uint) (map[uint][]model.EmployeeTenures, error) {
	var tenures []model.EmployeeTenures
	if err := r.Conn.Where("company_id = ? AND deleted_at IS NULL", companyId).
		Order("employee_id, joined_on DESC").Find(&tenures).Error; err != nil {
		return nil, err
	}
	result := make(map[uint][]model.EmployeeTenures)
	for _, t := range tenures {
		result[t.EmployeeId] = append(result[t.EmployeeId], t)
	}
	return result, nil
}

func (r *employeeRepository) GetAllDepartments(companyId uint) (map[uint][]model.Department, error) {
	type deptRow struct {
		ID         uint
		CompanyId  uint
		Name       string
		Depth      int
		OrderNo    int
		EmployeeId uint
	}
	var rows []deptRow
	if err := r.Conn.
		Table("departments").
		Select("departments.id, departments.company_id, departments.name, departments.depth, departments.order_no, employee_departments.employee_id").
		Joins("INNER JOIN employee_departments ON departments.id = employee_departments.department_id").
		Where("employee_departments.company_id = ? AND employee_departments.deleted_at IS NULL AND departments.deleted_at IS NULL", companyId).
		Scan(&rows).Error; err != nil {
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
