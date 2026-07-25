package usecase

import (
	"context"
	"mime/multipart"
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
