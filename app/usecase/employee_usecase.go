package usecase

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"app/domain/model"
	"app/domain/repository"
	"app/domain/service"
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
	ExportCSV(companyId uint) ([]byte, error)
	BulkCreateFromCSV(companyId uint, file multipart.File) error
	BulkUpdateFromCSV(companyId uint, file multipart.File) error
}

type employeeUseCase struct {
	csvService service.CsvService
	repo       repository.EmployeeRepository
	deptRepo   repository.DepartmentRepository
	prefRepo   repository.PrefectureRepository
}

func NewEmployeeUseCase(
	r repository.EmployeeRepository,
	deptRepo repository.DepartmentRepository,
	prefRepo repository.PrefectureRepository,
	csv service.CsvService,
) EmployeeUseCase {
	return &employeeUseCase{repo: r, deptRepo: deptRepo, prefRepo: prefRepo, csvService: csv}
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

func (u *employeeUseCase) ExportCSV(companyId uint) ([]byte, error) {
	employees, err := u.repo.GetList(companyId)
	if err != nil {
		return nil, err
	}
	addresses, err := u.repo.GetAllAddresses(companyId)
	if err != nil {
		return nil, err
	}
	tenures, err := u.repo.GetAllTenures(companyId)
	if err != nil {
		return nil, err
	}
	departments, err := u.repo.GetAllDepartments(companyId)
	if err != nil {
		return nil, err
	}

	exportData := make([]service.EmployeeExportData, len(employees))
	for i, emp := range employees {
		exportData[i] = service.EmployeeExportData{
			Employee:    emp,
			Address:     addresses[emp.ID],
			Tenures:     tenures[emp.ID],
			Departments: departments[emp.ID],
		}
	}
	return u.csvService.GenerateEmployeeCSV(exportData)
}

func (u *employeeUseCase) BulkCreateFromCSV(companyId uint, file multipart.File) error {
	rows, err := u.csvService.ParseEmployeeRows(file)
	if err != nil {
		return err
	}
	lookup, err := u.buildImportLookup(companyId)
	if err != nil {
		return err
	}

	// 既存スタッフコードが1件でもあればエラー
	var duplicates []string
	for _, row := range rows {
		if _, exists := lookup.empByCode[row.StaffCode]; exists {
			duplicates = append(duplicates, row.StaffCode)
		}
	}
	if len(duplicates) > 0 {
		return fmt.Errorf("以下のスタッフコードはすでに登録されています: %s", strings.Join(duplicates, ", "))
	}

	for _, row := range rows {
		emp := &model.Employee{
			CompanyId:     companyId,
			StaffCode:     row.StaffCode,
			LastName:      row.LastName,
			FirstName:     row.FirstName,
			LastNameKana:  row.LastNameKana,
			FirstNameKana: row.FirstNameKana,
			Email:         row.Email,
		}
		created, err := u.repo.Create(emp)
		if err != nil {
			return err
		}
		if err := u.applyRelated(companyId, created.ID, row, lookup); err != nil {
			return err
		}
	}
	return nil
}

func (u *employeeUseCase) BulkUpdateFromCSV(companyId uint, file multipart.File) error {
	rows, err := u.csvService.ParseEmployeeRows(file)
	if err != nil {
		return err
	}
	lookup, err := u.buildImportLookup(companyId)
	if err != nil {
		return err
	}

	// staff_code が未存在の行を先にすべてチェック
	var missing []string
	for _, row := range rows {
		if _, exists := lookup.empByCode[row.StaffCode]; !exists {
			missing = append(missing, row.StaffCode)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("以下のスタッフコードに紐づく社員が見つかりません: %s", strings.Join(missing, ", "))
	}

	for _, row := range rows {
		emp, exists := lookup.empByCode[row.StaffCode]
		if !exists {
			continue
		}
		emp.LastName = row.LastName
		emp.FirstName = row.FirstName
		emp.LastNameKana = row.LastNameKana
		emp.FirstNameKana = row.FirstNameKana
		emp.Email = row.Email
		if err := u.repo.UpdateEmployee(emp); err != nil {
			return err
		}
		if err := u.applyRelated(companyId, emp.ID, row, lookup); err != nil {
			return err
		}
	}
	return nil
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
	// 住所
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

	// 在籍情報（記載がある場合は全置換）
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

	// 所属部署（記載がある場合は全置換）
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

// gorm.ErrRecordNotFound を参照するためのブランクインポート抑止
var _ = gorm.ErrRecordNotFound
