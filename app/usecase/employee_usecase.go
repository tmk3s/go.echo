package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"app/domain/model"
	"app/domain/repository"
	"app/domain/service"
)

type TenureUpdateInput struct {
	ID              uint
	JoinedOn        string
	ResignationOn   *string
	ResignationType *string
	Status          *string
}

type UpdateAllInput struct {
	LastName      string
	FirstName     string
	LastNameKana  string
	FirstNameKana string
	Email         string
	StaffCode     string
	PostCode      *string
	PrefectureId  *uint
	City          *string
	AddressLine1  *string
	AddressLine2  *string
	Tel           *string
	Tenures       []TenureUpdateInput
	DepartmentIds []uint
}

// BulkImportEnqueuer is the interface for enqueuing bulk import background tasks.
type BulkImportEnqueuer interface {
	EnqueueBulkCreate(jobId uint) error
	EnqueueBulkUpdate(jobId uint) error
}

type EmployeeUseCase interface {
	GetEmployees(companyId uint) (*[]model.Employee, error)
	GetEmployeeDetail(companyId uint, id uint) (*model.Employee, error)
	Create(companyId uint, lastName string, firstName string, lastNameKana string, firstNameKana string, email string, staffCode string) error
	UpdateAll(companyId uint, id uint, input UpdateAllInput) error
	ExportCSV(companyId uint) ([]byte, error)
	EnqueueBulkCreate(companyId uint, file multipart.File) (*model.Job, error)
	EnqueueBulkUpdate(companyId uint, file multipart.File) (*model.Job, error)
	ProcessBulkCreate(ctx context.Context, jobId uint) error
	ProcessBulkUpdate(ctx context.Context, jobId uint) error
}

type employeeUseCase struct {
	csvService service.CsvService
	repo       repository.EmployeeRepository
	deptRepo   repository.DepartmentRepository
	prefRepo   repository.PrefectureRepository
	jobRepo    repository.JobRepository
	enqueuer   BulkImportEnqueuer
}

func NewEmployeeUseCase(
	r repository.EmployeeRepository,
	deptRepo repository.DepartmentRepository,
	prefRepo repository.PrefectureRepository,
	csv service.CsvService,
	jobRepo repository.JobRepository,
	enqueuer BulkImportEnqueuer,
) EmployeeUseCase {
	return &employeeUseCase{
		repo:       r,
		deptRepo:   deptRepo,
		prefRepo:   prefRepo,
		csvService: csv,
		jobRepo:    jobRepo,
		enqueuer:   enqueuer,
	}
}

func (u *employeeUseCase) GetEmployees(companyId uint) (*[]model.Employee, error) {
	employees, err := u.repo.GetList(companyId)
	if err != nil {
		return nil, err
	}
	return &employees, nil
}

func (u *employeeUseCase) GetEmployeeDetail(companyId uint, id uint) (*model.Employee, error) {
	return u.repo.GetDetail(companyId, id)
}

func (u *employeeUseCase) Create(companyId uint, lastName string, firstName string, lastNameKana string, firstNameKana string, email string, staffCode string) error {
	employee := &model.Employee{
		CompanyId:     companyId,
		LastName:      lastName,
		FirstName:     firstName,
		LastNameKana:  lastNameKana,
		FirstNameKana: firstNameKana,
		Email:         email,
		StaffCode:     staffCode,
	}
	_, err := u.repo.Create(employee)
	return err
}

func (u *employeeUseCase) UpdateAll(companyId uint, id uint, input UpdateAllInput) error {
	employee, err := u.repo.GetDetail(companyId, id)
	if err != nil {
		return err
	}
	employee.LastName = input.LastName
	employee.FirstName = input.FirstName
	employee.LastNameKana = input.LastNameKana
	employee.FirstNameKana = input.FirstNameKana
	employee.Email = input.Email
	employee.StaffCode = input.StaffCode
	if err := u.repo.UpdateEmployee(employee); err != nil {
		return err
	}

	address := &model.EmployeeAddress{
		CompanyId:    companyId,
		EmployeeId:   id,
		PostCode:     input.PostCode,
		PrefectureId: input.PrefectureId,
		City:         input.City,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		Tel:          input.Tel,
	}
	if err := u.repo.UpsertAddress(address); err != nil {
		return err
	}

	for _, t := range input.Tenures {
		target := findTenure(employee.Tenures, t.ID)
		if target == nil {
			continue
		}
		parsed, err := time.Parse("2006-01-02", t.JoinedOn)
		if err != nil {
			return err
		}
		target.JoinedOn = parsed
		if t.ResignationOn != nil && *t.ResignationOn != "" {
			rt, err := time.Parse("2006-01-02", *t.ResignationOn)
			if err != nil {
				return err
			}
			target.ResignationOn = &rt
		} else {
			target.ResignationOn = nil
		}
		target.ResignationType = t.ResignationType
		target.Status = t.Status
		if err := u.repo.UpdateTenure(target); err != nil {
			return err
		}
	}

	return u.repo.UpdateDepartments(companyId, id, input.DepartmentIds)
}

func findTenure(tenures []model.EmployeeTenures, id uint) *model.EmployeeTenures {
	for i := range tenures {
		if tenures[i].ID == id {
			return &tenures[i]
		}
	}
	return nil
}

func (u *employeeUseCase) ExportCSV(companyId uint) ([]byte, error) {
	employees, err := u.repo.GetListForExport(companyId)
	if err != nil {
		return nil, err
	}
	return u.csvService.GenerateEmployeeCSV(employees)
}

func (u *employeeUseCase) EnqueueBulkCreate(companyId uint, file multipart.File) (*model.Job, error) {
	csvBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	job, err := u.jobRepo.Create(&model.Job{
		CompanyId: companyId,
		JobType:   "employee:bulk_create",
		Status:    "pending",
		FileData:  csvBytes,
	})
	if err != nil {
		return nil, err
	}

	if err := u.enqueuer.EnqueueBulkCreate(job.ID); err != nil {
		u.jobRepo.Fail(job.ID, err.Error())
		return nil, err
	}
	return job, nil
}

func (u *employeeUseCase) EnqueueBulkUpdate(companyId uint, file multipart.File) (*model.Job, error) {
	csvBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	job, err := u.jobRepo.Create(&model.Job{
		CompanyId: companyId,
		JobType:   "employee:bulk_update",
		Status:    "pending",
		FileData:  csvBytes,
	})
	if err != nil {
		return nil, err
	}

	if err := u.enqueuer.EnqueueBulkUpdate(job.ID); err != nil {
		u.jobRepo.Fail(job.ID, err.Error())
		return nil, err
	}
	return job, nil
}

// csvReadCloser wraps bytes.Reader to satisfy multipart.File (adds no-op Close).
type csvReadCloser struct{ *bytes.Reader }

func (csvReadCloser) Close() error { return nil }

func (u *employeeUseCase) ProcessBulkCreate(ctx context.Context, jobId uint) error {
	job, err := u.jobRepo.GetById(jobId)
	if err != nil {
		return err
	}

	rows, err := u.csvService.ParseEmployeeRows(csvReadCloser{bytes.NewReader(job.FileData)})
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	lookup, err := u.buildImportLookup(job.CompanyId)
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	var duplicates []string
	for _, row := range rows {
		if _, exists := lookup.empByCode[row.StaffCode]; exists {
			duplicates = append(duplicates, row.StaffCode)
		}
	}
	if len(duplicates) > 0 {
		msg := fmt.Sprintf("以下のスタッフコードはすでに登録されています: %s", strings.Join(duplicates, ", "))
		u.jobRepo.Fail(jobId, msg)
		return fmt.Errorf("%s", msg)
	}

	if err := u.jobRepo.StartProcessing(jobId, len(rows)); err != nil {
		return err
	}

	for i, row := range rows {
		select {
		case <-ctx.Done():
			u.jobRepo.Fail(jobId, "処理がキャンセルされました")
			return ctx.Err()
		default:
		}

		emp := &model.Employee{
			CompanyId:     job.CompanyId,
			StaffCode:     row.StaffCode,
			LastName:      row.LastName,
			FirstName:     row.FirstName,
			LastNameKana:  row.LastNameKana,
			FirstNameKana: row.FirstNameKana,
			Email:         row.Email,
		}
		created, err := u.repo.Create(emp)
		if err != nil {
			u.jobRepo.Fail(jobId, fmt.Sprintf("行%d: %v", i+1, err))
			return err
		}
		if err := u.applyRelated(job.CompanyId, created.ID, row, lookup); err != nil {
			u.jobRepo.Fail(jobId, fmt.Sprintf("行%d: %v", i+1, err))
			return err
		}
		u.jobRepo.UpdateProgress(jobId, i+1)
	}

	return u.jobRepo.Complete(jobId)
}

func (u *employeeUseCase) ProcessBulkUpdate(ctx context.Context, jobId uint) error {
	job, err := u.jobRepo.GetById(jobId)
	if err != nil {
		return err
	}

	rows, err := u.csvService.ParseEmployeeRows(csvReadCloser{bytes.NewReader(job.FileData)})
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	lookup, err := u.buildImportLookup(job.CompanyId)
	if err != nil {
		u.jobRepo.Fail(jobId, err.Error())
		return err
	}

	var missing []string
	for _, row := range rows {
		if _, exists := lookup.empByCode[row.StaffCode]; !exists {
			missing = append(missing, row.StaffCode)
		}
	}
	if len(missing) > 0 {
		msg := fmt.Sprintf("以下のスタッフコードに紐づく社員が見つかりません: %s", strings.Join(missing, ", "))
		u.jobRepo.Fail(jobId, msg)
		return fmt.Errorf("%s", msg)
	}

	if err := u.jobRepo.StartProcessing(jobId, len(rows)); err != nil {
		return err
	}

	for i, row := range rows {
		select {
		case <-ctx.Done():
			u.jobRepo.Fail(jobId, "処理がキャンセルされました")
			return ctx.Err()
		default:
		}

		emp := lookup.empByCode[row.StaffCode]
		emp.LastName = row.LastName
		emp.FirstName = row.FirstName
		emp.LastNameKana = row.LastNameKana
		emp.FirstNameKana = row.FirstNameKana
		emp.Email = row.Email
		if err := u.repo.UpdateEmployee(emp); err != nil {
			u.jobRepo.Fail(jobId, fmt.Sprintf("行%d: %v", i+1, err))
			return err
		}
		if err := u.applyRelated(job.CompanyId, emp.ID, row, lookup); err != nil {
			u.jobRepo.Fail(jobId, fmt.Sprintf("行%d: %v", i+1, err))
			return err
		}
		u.jobRepo.UpdateProgress(jobId, i+1)
	}

	return u.jobRepo.Complete(jobId)
}

type importLookup struct {
	empByCode  map[string]*model.Employee
	deptByName map[string]uint
	prefByName map[string]uint
}

func (u *employeeUseCase) buildImportLookup(companyId uint) (*importLookup, error) {
	emps, err := u.repo.GetList(companyId)
	if err != nil {
		return nil, err
	}
	empByCode := make(map[string]*model.Employee, len(emps))
	for i := range emps {
		empByCode[emps[i].StaffCode] = &emps[i]
	}

	depts, err := u.deptRepo.GetList(companyId)
	if err != nil {
		return nil, err
	}
	deptByName := make(map[string]uint, len(depts))
	for _, d := range depts {
		deptByName[d.Name] = d.ID
	}

	prefs, err := u.prefRepo.GetAll()
	if err != nil {
		return nil, err
	}
	prefByName := make(map[string]uint, len(prefs))
	for _, p := range prefs {
		prefByName[p.Name] = p.ID
	}

	return &importLookup{empByCode: empByCode, deptByName: deptByName, prefByName: prefByName}, nil
}

func (u *employeeUseCase) applyRelated(companyId uint, empID uint, row service.EmployeeCSVRow, l *importLookup) error {
	if row.PostCode != "" || row.PrefectureName != "" || row.City != "" ||
		row.AddressLine1 != "" || row.AddressLine2 != "" || row.Tel != "" {
		var prefId *uint
		if row.PrefectureName != "" {
			if id, ok := l.prefByName[row.PrefectureName]; ok {
				prefId = &id
			}
		}
		addr := &model.EmployeeAddress{
			CompanyId:    companyId,
			EmployeeId:   empID,
			PostCode:     nilStr(row.PostCode),
			PrefectureId: prefId,
			City:         nilStr(row.City),
			AddressLine1: nilStr(row.AddressLine1),
			AddressLine2: nilStr(row.AddressLine2),
			Tel:          nilStr(row.Tel),
		}
		if err := u.repo.UpsertAddress(addr); err != nil {
			return err
		}
	}

	if len(row.Tenures) > 0 {
		tenures := make([]model.EmployeeTenures, 0, len(row.Tenures))
		for _, t := range row.Tenures {
			parsed, err := time.Parse("2006-01-02", t.JoinedOn)
			if err != nil {
				return err
			}
			tenure := model.EmployeeTenures{CompanyId: companyId, EmployeeId: empID, JoinedOn: parsed}
			if t.ResignationOn != "" {
				rt, err := time.Parse("2006-01-02", t.ResignationOn)
				if err != nil {
					return err
				}
				tenure.ResignationOn = &rt
			}
			if t.ResignationType != "" {
				tenure.ResignationType = &t.ResignationType
			}
			if t.Status != "" {
				tenure.Status = &t.Status
			}
			tenures = append(tenures, tenure)
		}
		if err := u.repo.ReplaceTenures(companyId, empID, tenures); err != nil {
			return err
		}
	}

	if len(row.DepartmentNames) > 0 {
		var deptIds []uint
		for _, name := range row.DepartmentNames {
			if id, ok := l.deptByName[name]; ok {
				deptIds = append(deptIds, id)
			}
		}
		if err := u.repo.UpdateDepartments(companyId, empID, deptIds); err != nil {
			return err
		}
	}
	return nil
}

func nilStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
