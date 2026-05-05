package usecase

import (
	"time"

	"app/domain/model"
	"app/domain/repository"
	"gorm.io/gorm"
)

type EmployeeDetail struct {
	Employee    *model.Employee
	Address     *model.EmployeeAddress
	Tenures     []model.EmployeeTenures
	Departments []model.Department
}

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
	GetEmployeeDetail(companyId uint, id uint) (*EmployeeDetail, error)
	Create(companyId uint, lastName string, firstName string, lastNameKana string, firstNameKana string, email string, staffCode string) error
	UpdateAll(companyId uint, id uint, input UpdateAllInput) error
}

type employeeUseCase struct {
	repo repository.EmployeeRepository
}

func NewEmployeeUseCase(r repository.EmployeeRepository) EmployeeUseCase {
	return &employeeUseCase{repo: r}
}

func (u *employeeUseCase) GetEmployees(companyId uint) (*[]model.Employee, error) {
	employees, err := u.repo.GetList(companyId)
	if err != nil {
		return nil, err
	}
	return &employees, nil
}

func (u *employeeUseCase) GetEmployeeDetail(companyId uint, id uint) (*EmployeeDetail, error) {
	employee, err := u.repo.GetById(companyId, id)
	if err != nil {
		return nil, err
	}
	address, err := u.repo.GetAddress(companyId, id)
	if err != nil {
		return nil, err
	}
	tenures, err := u.repo.GetTenures(companyId, id)
	if err != nil {
		return nil, err
	}
	departments, err := u.repo.GetDepartments(companyId, id)
	if err != nil {
		return nil, err
	}
	return &EmployeeDetail{
		Employee:    employee,
		Address:     address,
		Tenures:     tenures,
		Departments: departments,
	}, nil
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
	employee, err := u.repo.GetById(companyId, id)
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

	existingTenures, err := u.repo.GetTenures(companyId, id)
	if err != nil {
		return err
	}
	tenureMap := make(map[uint]*model.EmployeeTenures, len(existingTenures))
	for i := range existingTenures {
		tenureMap[existingTenures[i].ID] = &existingTenures[i]
	}
	for _, t := range input.Tenures {
		target, ok := tenureMap[t.ID]
		if !ok {
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

// gorm.ErrRecordNotFound を参照するためのブランクインポート抑止
var _ = gorm.ErrRecordNotFound
