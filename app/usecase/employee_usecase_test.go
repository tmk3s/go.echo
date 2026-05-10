package usecase_test

import (
	"context"
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
	return usecase.NewEmployeeUseCase(repo, dept, pref, csv, &mockJobRepo{}, &mockBulkImportEnqueuer{})
}

func newUCWithJob(
	repo *mockEmployeeRepo,
	dept *mockDeptRepo,
	pref *mockPrefRepo,
	csv *mockCsvService,
	jobRepo *mockJobRepo,
	enqueuer *mockBulkImportEnqueuer,
) usecase.EmployeeUseCase {
	return usecase.NewEmployeeUseCase(repo, dept, pref, csv, jobRepo, enqueuer)
}

// seedJob creates a job in the repo and returns its ID.
func seedJob(t *testing.T, jobRepo *mockJobRepo, companyId uint) uint {
	t.Helper()
	job, err := jobRepo.Create(&model.Job{CompanyId: companyId, JobType: "employee:bulk_create", Status: "pending"})
	if err != nil {
		t.Fatalf("seedJob: %v", err)
	}
	return job.ID
}

// ---- EnqueueBulkCreate ----

func TestEnqueueBulkCreate_SavesCSVAndEnqueues(t *testing.T) {
	jobRepo := &mockJobRepo{}
	enqueuer := &mockBulkImportEnqueuer{}

	job, err := newUCWithJob(&mockEmployeeRepo{}, &mockDeptRepo{}, &mockPrefRepo{}, &mockCsvService{}, jobRepo, enqueuer).
		EnqueueBulkCreate(1, emptyFile())
	if err != nil {
		t.Fatal(err)
	}
	if job == nil || job.ID == 0 {
		t.Fatal("expected job to be created with ID")
	}
	if len(enqueuer.enqueuedCreateJobIds) != 1 || enqueuer.enqueuedCreateJobIds[0] != job.ID {
		t.Errorf("enqueuer not called with correct job ID, got %v", enqueuer.enqueuedCreateJobIds)
	}
	if job.JobType != "employee:bulk_create" {
		t.Errorf("want job_type employee:bulk_create, got %s", job.JobType)
	}
}

// ---- ProcessBulkCreate ----

func TestProcessBulkCreate_DuplicateStaffCode(t *testing.T) {
	repo := &mockEmployeeRepo{
		employees: []model.Employee{{StaffCode: "EMP001"}},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP001", Email: "emp@example.com"},
		},
	}
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	err := uc.ProcessBulkCreate(context.Background(), jobId)
	if err == nil {
		t.Fatal("want error for duplicate staff code, got nil")
	}
	if !strings.Contains(err.Error(), "EMP001") {
		t.Errorf("error should mention EMP001, got: %s", err.Error())
	}
	if jobRepo.finalStatus != "failed" {
		t.Errorf("want job status failed, got %s", jobRepo.finalStatus)
	}
}

func TestProcessBulkCreate_MultipleDuplicates(t *testing.T) {
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
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	err := uc.ProcessBulkCreate(context.Background(), jobId)
	if err == nil {
		t.Fatal("want error for duplicate staff codes, got nil")
	}
	if !strings.Contains(err.Error(), "EMP001") || !strings.Contains(err.Error(), "EMP002") {
		t.Errorf("error should list both duplicates, got: %s", err.Error())
	}
}

func TestProcessBulkCreate_CreatesEmployee(t *testing.T) {
	var createdEmp *model.Employee
	repo := &mockEmployeeRepo{
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
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	if err := uc.ProcessBulkCreate(context.Background(), jobId); err != nil {
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
	if jobRepo.finalStatus != "completed" {
		t.Errorf("want status completed, got %s", jobRepo.finalStatus)
	}
}

func TestProcessBulkCreate_WithDepartment(t *testing.T) {
	var assignedDeptIds []uint
	repo := &mockEmployeeRepo{
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
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, dept, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	if err := uc.ProcessBulkCreate(context.Background(), jobId); err != nil {
		t.Fatal(err)
	}
	if len(assignedDeptIds) != 1 || assignedDeptIds[0] != 10 {
		t.Errorf("UpdateDepartments called with %v, want [10]", assignedDeptIds)
	}
}

func TestProcessBulkCreate_UpdatesProgress(t *testing.T) {
	repo := &mockEmployeeRepo{
		createFunc: func(emp *model.Employee) (*model.Employee, error) {
			emp.ID = 1
			return emp, nil
		},
	}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "E001", Email: "e1@example.com"},
			{StaffCode: "E002", Email: "e2@example.com"},
		},
	}
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	if err := uc.ProcessBulkCreate(context.Background(), jobId); err != nil {
		t.Fatal(err)
	}
	if len(jobRepo.progressLog) != 2 {
		t.Errorf("want 2 progress updates, got %d", len(jobRepo.progressLog))
	}
	if jobRepo.progressLog[1] != 2 {
		t.Errorf("want final progress=2, got %d", jobRepo.progressLog[1])
	}
}

// ---- EnqueueBulkUpdate ----

func TestEnqueueBulkUpdate_SavesCSVAndEnqueues(t *testing.T) {
	jobRepo := &mockJobRepo{}
	enqueuer := &mockBulkImportEnqueuer{}

	job, err := newUCWithJob(&mockEmployeeRepo{}, &mockDeptRepo{}, &mockPrefRepo{}, &mockCsvService{}, jobRepo, enqueuer).
		EnqueueBulkUpdate(1, emptyFile())
	if err != nil {
		t.Fatal(err)
	}
	if job == nil || job.ID == 0 {
		t.Fatal("expected job to be created with ID")
	}
	if len(enqueuer.enqueuedUpdateJobIds) != 1 || enqueuer.enqueuedUpdateJobIds[0] != job.ID {
		t.Errorf("enqueuer not called with correct job ID, got %v", enqueuer.enqueuedUpdateJobIds)
	}
	if job.JobType != "employee:bulk_update" {
		t.Errorf("want job_type employee:bulk_update, got %s", job.JobType)
	}
}

// ---- ProcessBulkUpdate ----

func TestProcessBulkUpdate_MissingStaffCode(t *testing.T) {
	repo := &mockEmployeeRepo{employees: []model.Employee{}}
	csv := &mockCsvService{
		rows: []domainservice.EmployeeCSVRow{
			{StaffCode: "EMP999", Email: "emp@example.com"},
		},
	}
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	err := uc.ProcessBulkUpdate(context.Background(), jobId)
	if err == nil {
		t.Fatal("want error for missing staff code, got nil")
	}
	if !strings.Contains(err.Error(), "EMP999") {
		t.Errorf("error should mention EMP999, got: %s", err.Error())
	}
	if jobRepo.finalStatus != "failed" {
		t.Errorf("want job status failed, got %s", jobRepo.finalStatus)
	}
}

func TestProcessBulkUpdate_MultipleMissing(t *testing.T) {
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
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	err := uc.ProcessBulkUpdate(context.Background(), jobId)
	if err == nil {
		t.Fatal("want error for missing staff codes, got nil")
	}
	if !strings.Contains(err.Error(), "EMP888") || !strings.Contains(err.Error(), "EMP999") {
		t.Errorf("error should list missing codes, got: %s", err.Error())
	}
}

func TestProcessBulkUpdate_UpdatesEmployee(t *testing.T) {
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
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	if err := uc.ProcessBulkUpdate(context.Background(), jobId); err != nil {
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

func TestProcessBulkUpdate_AllCodes(t *testing.T) {
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
	jobRepo := &mockJobRepo{}
	jobId := seedJob(t, jobRepo, 1)

	uc := newUCWithJob(repo, &mockDeptRepo{}, &mockPrefRepo{}, csv, jobRepo, &mockBulkImportEnqueuer{})
	if err := uc.ProcessBulkUpdate(context.Background(), jobId); err != nil {
		t.Fatal(err)
	}
	if updateCount != 2 {
		t.Errorf("want UpdateEmployee called 2 times, got %d", updateCount)
	}
}
