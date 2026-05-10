package usecase_test

import (
	"strings"
	"testing"

	"app/domain/model"
	domainservice "app/domain/service"
	"app/usecase"
)

func newUC(
	repo *mockEmployeeRepo,
	dept *mockDeptRepo,
	pref *mockPrefRepo,
	csv *mockCsvService,
) usecase.EmployeeUseCase {
	return usecase.NewEmployeeUseCase(repo, dept, pref, csv)
}

// ---- BulkCreateFromCSV ----

func TestBulkCreateFromCSV_DuplicateStaffCode(t *testing.T) {
	repo := &mockEmployeeRepo{
		employees: []model.Employee{{StaffCode: "EMP001"}},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP001", Email: "emp@example.com"},
		},
	}

	err := newUC(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv).BulkCreateFromCSV(1, emptyFile())
	if err == nil {
		t.Fatal("want error for duplicate staff code, got nil")
	}
	if !strings.Contains(err.Error(), "EMP001") {
		t.Errorf("error should mention EMP001, got: %s", err.Error())
	}
}

func TestBulkCreateFromCSV_MultipleDuplicates(t *testing.T) {
	repo := &mockEmployeeRepo{
		employees: []model.Employee{
			{StaffCode: "EMP001"},
			{StaffCode: "EMP002"},
		},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP001", Email: "e1@example.com"},
			{StaffCode: "EMP002", Email: "e2@example.com"},
			{StaffCode: "EMP003", Email: "e3@example.com"},
		},
	}

	err := newUC(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv).BulkCreateFromCSV(1, emptyFile())
	if err == nil {
		t.Fatal("want error for duplicate staff codes, got nil")
	}
	if !strings.Contains(err.Error(), "EMP001") || !strings.Contains(err.Error(), "EMP002") {
		t.Errorf("error should list both duplicates, got: %s", err.Error())
	}
}

func TestBulkCreateFromCSV_Success(t *testing.T) {
	var createdEmp *model.Employee
	repo := &mockEmployeeRepo{
		employees: []model.Employee{},
		createFunc: func(emp *model.Employee) (*model.Employee, error) {
			emp.ID = 99
			createdEmp = emp
			return emp, nil
		},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP099", LastName: "山田", Email: "yamada@example.com"},
		},
	}

	if err := newUC(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv).BulkCreateFromCSV(1, emptyFile()); err != nil {
		t.Fatal(err)
	}
	if createdEmp == nil {
		t.Fatal("Create was not called")
	}
	if createdEmp.StaffCode != "EMP099" {
		t.Errorf("StaffCode: want EMP099, got %s", createdEmp.StaffCode)
	}
	if createdEmp.LastName != "山田" {
		t.Errorf("LastName: want 山田, got %s", createdEmp.LastName)
	}
}

func TestBulkCreateFromCSV_WithDepartment(t *testing.T) {
	var assignedDeptIds []uint
	repo := &mockEmployeeRepo{
		employees: []model.Employee{},
		createFunc: func(emp *model.Employee) (*model.Employee, error) {
			emp.ID = 1
			return emp, nil
		},
		updateDepts: func(_, _ uint, ids []uint) error {
			assignedDeptIds = ids
			return nil
		},
	}
	dept := &mockDeptRepo{
		departments: []model.Department{{Name: "営業部"}},
	}
	dept.departments[0].ID = 10

	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP001", Email: "e@example.com", DepartmentNames: []string{"営業部"}},
		},
	}

	if err := newUC(repo, dept, &mockPrefRepo{}, csv).BulkCreateFromCSV(1, emptyFile()); err != nil {
		t.Fatal(err)
	}
	if len(assignedDeptIds) != 1 || assignedDeptIds[0] != 10 {
		t.Errorf("UpdateDepartments called with %v, want [10]", assignedDeptIds)
	}
}

// ---- BulkUpdateFromCSV ----

func TestBulkUpdateFromCSV_MissingStaffCode(t *testing.T) {
	repo := &mockEmployeeRepo{employees: []model.Employee{}}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP999", Email: "emp@example.com"},
		},
	}

	err := newUC(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv).BulkUpdateFromCSV(1, emptyFile())
	if err == nil {
		t.Fatal("want error for missing staff code, got nil")
	}
	if !strings.Contains(err.Error(), "EMP999") {
		t.Errorf("error should mention EMP999, got: %s", err.Error())
	}
}

func TestBulkUpdateFromCSV_MultipleMissing(t *testing.T) {
	repo := &mockEmployeeRepo{
		employees: []model.Employee{{StaffCode: "EMP001"}},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP001", Email: "e1@example.com"},
			{StaffCode: "EMP888", Email: "e888@example.com"},
			{StaffCode: "EMP999", Email: "e999@example.com"},
		},
	}

	err := newUC(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv).BulkUpdateFromCSV(1, emptyFile())
	if err == nil {
		t.Fatal("want error for missing staff codes, got nil")
	}
	if !strings.Contains(err.Error(), "EMP888") || !strings.Contains(err.Error(), "EMP999") {
		t.Errorf("error should list missing codes, got: %s", err.Error())
	}
}

func TestBulkUpdateFromCSV_Success(t *testing.T) {
	var updatedEmp *model.Employee
	existing := model.Employee{StaffCode: "EMP001", LastName: "旧姓", CompanyId: 1}
	existing.ID = 1

	repo := &mockEmployeeRepo{
		employees: []model.Employee{existing},
		updateFunc: func(emp *model.Employee) error {
			updatedEmp = emp
			return nil
		},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP001", LastName: "新姓", FirstName: "花子", Email: "e@example.com"},
		},
	}

	if err := newUC(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv).BulkUpdateFromCSV(1, emptyFile()); err != nil {
		t.Fatal(err)
	}
	if updatedEmp == nil {
		t.Fatal("UpdateEmployee was not called")
	}
	if updatedEmp.LastName != "新姓" {
		t.Errorf("LastName: want 新姓, got %s", updatedEmp.LastName)
	}
	if updatedEmp.FirstName != "花子" {
		t.Errorf("FirstName: want 花子, got %s", updatedEmp.FirstName)
	}
}

func TestBulkUpdateFromCSV_AllExistingCodes(t *testing.T) {
	var updateCount int
	employees := []model.Employee{
		{StaffCode: "EMP001"},
		{StaffCode: "EMP002"},
	}
	employees[0].ID = 1
	employees[1].ID = 2

	repo := &mockEmployeeRepo{
		employees: employees,
		updateFunc: func(_ *model.Employee) error {
			updateCount++
			return nil
		},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP001", Email: "e1@example.com"},
			{StaffCode: "EMP002", Email: "e2@example.com"},
		},
	}

	if err := newUC(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv).BulkUpdateFromCSV(1, emptyFile()); err != nil {
		t.Fatal(err)
	}
	if updateCount != 2 {
		t.Errorf("want UpdateEmployee called 2 times, got %d", updateCount)
	}
}
