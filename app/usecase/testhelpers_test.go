package usecase_test

import (
	"bytes"
	"errors"
	"mime/multipart"

	"app/domain/model"
	domainservice "app/domain/service"
)

var errNotFound = errors.New("record not found")

// ---- file helper ----

type nopFile struct{ *bytes.Reader }

func (nopFile) Close() error { return nil }

func emptyFile() multipart.File { return nopFile{bytes.NewReader(nil)} }

// ---- mockUserRepo ----

type mockUserRepo struct {
	user      *model.User
	findErr   error
	createErr error
}

func (m *mockUserRepo) GetById(_ uint) (*model.User, error)                { return m.user, m.findErr }
func (m *mockUserRepo) GetByEmail(_ string) (*model.User, error)           { return m.user, m.findErr }
func (m *mockUserRepo) GetByEmailAndPass(_, _ string) (*model.User, error) { return m.user, m.findErr }
func (m *mockUserRepo) Create(u *model.User) (*model.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	u.ID = 1
	return u, nil
}
func (m *mockUserRepo) Update(u *model.User) (*model.User, error) { return u, nil }
func (m *mockUserRepo) Delete(_ uint) error                       { return nil }

// ---- mockTodoRepo ----

type mockTodoRepo struct {
	todo         *model.Todo
	todos        []model.Todo
	getErr       error
	addCalled    bool
	updatedTodo  *model.Todo
	deleteCalled bool
}

func (m *mockTodoRepo) GetById(_ uint) (*model.Todo, error) { return m.todo, m.getErr }
func (m *mockTodoRepo) GetList(_ uint) ([]model.Todo, error) { return m.todos, nil }
func (m *mockTodoRepo) Add(todo *model.Todo) (*model.Todo, error) {
	m.addCalled = true
	todo.ID = 1
	return todo, nil
}
func (m *mockTodoRepo) Update(todo *model.Todo) (*model.Todo, error) {
	m.updatedTodo = todo
	return todo, nil
}
func (m *mockTodoRepo) Delete(_ *model.Todo) error {
	m.deleteCalled = true
	return nil
}

// ---- mockEmployeeRepo ----

type mockEmployeeRepo struct {
	employees      []model.Employee
	createFunc     func(*model.Employee) (*model.Employee, error)
	updateFunc     func(*model.Employee) error
	upsertFunc     func(*model.EmployeeAddress) error
	replaceTenures func(uint, uint, []model.EmployeeTenures) error
	updateDepts    func(uint, uint, []uint) error
}

func (m *mockEmployeeRepo) GetList(_ uint) ([]model.Employee, error) { return m.employees, nil }
func (m *mockEmployeeRepo) GetDetail(_, _ uint) (*model.Employee, error) { return nil, nil }
func (m *mockEmployeeRepo) GetListForExport(_ uint) ([]model.Employee, error) { return nil, nil }
func (m *mockEmployeeRepo) Create(emp *model.Employee) (*model.Employee, error) {
	if m.createFunc != nil {
		return m.createFunc(emp)
	}
	emp.ID = 1
	return emp, nil
}
func (m *mockEmployeeRepo) BulkCreate(_ []model.Employee) error { return nil }
func (m *mockEmployeeRepo) UpdateEmployee(emp *model.Employee) error {
	if m.updateFunc != nil {
		return m.updateFunc(emp)
	}
	return nil
}
func (m *mockEmployeeRepo) UpsertAddress(addr *model.EmployeeAddress) error {
	if m.upsertFunc != nil {
		return m.upsertFunc(addr)
	}
	return nil
}
func (m *mockEmployeeRepo) UpdateTenure(_ *model.EmployeeTenures) error { return nil }
func (m *mockEmployeeRepo) UpdateDepartments(cid, eid uint, dids []uint) error {
	if m.updateDepts != nil {
		return m.updateDepts(cid, eid, dids)
	}
	return nil
}
func (m *mockEmployeeRepo) ReplaceTenures(cid, eid uint, tenures []model.EmployeeTenures) error {
	if m.replaceTenures != nil {
		return m.replaceTenures(cid, eid, tenures)
	}
	return nil
}

// ---- mockDeptRepo ----

type mockDeptRepo struct {
	byId        *model.Department
	departments []model.Department
	createCount int
}

func (m *mockDeptRepo) GetById(_ uint) (*model.Department, error) {
	if m.byId != nil {
		return m.byId, nil
	}
	return nil, errNotFound
}
func (m *mockDeptRepo) GetList(_ uint) ([]model.Department, error) { return m.departments, nil }
func (m *mockDeptRepo) Create(d *model.Department, _ *uint) (*model.Department, error) {
	m.createCount++
	d.ID = uint(m.createCount)
	return d, nil
}
func (m *mockDeptRepo) Update(d *model.Department) (*model.Department, error) { return d, nil }
func (m *mockDeptRepo) Delete(_ *model.Department) error                      { return nil }

// ---- mockPrefRepo ----

type mockPrefRepo struct {
	prefectures []model.Prefecture
}

func (m *mockPrefRepo) GetAll() ([]model.Prefecture, error) { return m.prefectures, nil }

// ---- mockCompanyRepo ----

type mockCompanyRepo struct {
	company     *model.Company
	updatedName string
}

func (m *mockCompanyRepo) GetById(_ uint) (*model.Company, error) { return m.company, nil }
func (m *mockCompanyRepo) Update(c *model.Company) error {
	m.updatedName = c.Name
	return nil
}

// ---- mockCsvService ----

type mockCsvService struct {
	rows      []domainservice.EmployeeCSVRow
	deptNames []string
	err       error
}

func (m *mockCsvService) ParseDepartmentNames(_ multipart.File) ([]string, error) {
	return m.deptNames, m.err
}
func (m *mockCsvService) GenerateEmployeeCSV(_ []model.Employee) ([]byte, error) { return nil, nil }
func (m *mockCsvService) ParseEmployeeRows(_ multipart.File) ([]domainservice.EmployeeCSVRow, error) {
	return m.rows, m.err
}
