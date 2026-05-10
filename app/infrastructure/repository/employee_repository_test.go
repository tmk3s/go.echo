package repository_test

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"app/domain/model"
	"app/infrastructure/repository"
)

func setupEmployeeDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := db.AutoMigrate(
		&model.Employee{},
		&model.EmployeeAddress{},
		&model.EmployeeTenures{},
		&model.EmployeeDepartments{},
		&model.Department{},
		&model.Prefecture{},
	); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		truncate(t, db, "employee_departments", "employee_tenures", "employee_addresses", "employees", "departments", "prefectures")
	})
	return db
}

func seedEmployee(t *testing.T, db *gorm.DB, companyId uint, staffCode string) model.Employee {
	t.Helper()
	emp := model.Employee{CompanyId: companyId, StaffCode: staffCode, Email: staffCode + "@example.com"}
	if err := db.Create(&emp).Error; err != nil {
		t.Fatal(err)
	}
	return emp
}

func parseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

// --- Create ---

func TestEmployeeCreate_AssignsID(t *testing.T) {
	repo := repository.NewEmployeeRepository(setupEmployeeDB(t))
	emp, err := repo.Create(&model.Employee{
		CompanyId: 1, StaffCode: "EMP001", Email: "emp@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if emp.ID == 0 {
		t.Error("expected ID to be assigned after Create")
	}
}

func TestEmployeeCreate_StoresCorrectFields(t *testing.T) {
	repo := repository.NewEmployeeRepository(setupEmployeeDB(t))
	emp, err := repo.Create(&model.Employee{
		CompanyId:     1,
		StaffCode:     "EMP001",
		LastName:      "田中",
		FirstName:     "太郎",
		LastNameKana:  "タナカ",
		FirstNameKana: "タロウ",
		Email:         "taro@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if emp.LastName != "田中" {
		t.Errorf("LastName: want 田中, got %s", emp.LastName)
	}
	if emp.StaffCode != "EMP001" {
		t.Errorf("StaffCode: want EMP001, got %s", emp.StaffCode)
	}
}

// --- GetList ---

func TestEmployeeGetList_FiltersByCompany(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)

	seedEmployee(t, db, 1, "E001")
	seedEmployee(t, db, 1, "E002")
	seedEmployee(t, db, 2, "E003")

	list, err := repo.GetList(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("want 2 employees for company 1, got %d", len(list))
	}
	for _, e := range list {
		if e.CompanyId != 1 {
			t.Errorf("got employee with CompanyId=%d, want 1", e.CompanyId)
		}
	}
}

func TestEmployeeGetList_EmptyForUnknownCompany(t *testing.T) {
	repo := repository.NewEmployeeRepository(setupEmployeeDB(t))

	list, err := repo.GetList(999)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Errorf("want 0 employees, got %d", len(list))
	}
}

// --- GetDetail ---

func TestEmployeeGetDetail_ReturnsEmployee(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)
	emp := seedEmployee(t, db, 1, "EMP001")

	got, err := repo.GetDetail(1, emp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.StaffCode != "EMP001" {
		t.Errorf("want EMP001, got %s", got.StaffCode)
	}
}

func TestEmployeeGetDetail_PreloadsAddress(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)
	emp := seedEmployee(t, db, 1, "EMP001")

	postCode := "100-0001"
	db.Create(&model.EmployeeAddress{CompanyId: 1, EmployeeId: emp.ID, PostCode: &postCode})

	got, err := repo.GetDetail(1, emp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Address == nil {
		t.Fatal("expected address to be preloaded, got nil")
	}
	if got.Address.PostCode == nil || *got.Address.PostCode != "100-0001" {
		t.Errorf("PostCode: want 100-0001, got %v", got.Address.PostCode)
	}
}

func TestEmployeeGetDetail_PreloadsTenures(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)
	emp := seedEmployee(t, db, 1, "EMP001")

	status := "在職中"
	db.Create(&model.EmployeeTenures{
		CompanyId:  1,
		EmployeeId: emp.ID,
		JoinedOn:   parseDate("2020-04-01"),
		Status:     &status,
	})

	got, err := repo.GetDetail(1, emp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tenures) != 1 {
		t.Fatalf("want 1 tenure, got %d", len(got.Tenures))
	}
	if got.Tenures[0].Status == nil || *got.Tenures[0].Status != "在職中" {
		t.Errorf("Status: want 在職中, got %v", got.Tenures[0].Status)
	}
}

// --- UpsertAddress ---

func TestUpsertAddress_CreatesWhenNotExists(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)
	emp := seedEmployee(t, db, 1, "EA001")

	postCode := "100-0001"
	if err := repo.UpsertAddress(&model.EmployeeAddress{
		CompanyId: 1, EmployeeId: emp.ID, PostCode: &postCode,
	}); err != nil {
		t.Fatal(err)
	}

	var addr model.EmployeeAddress
	if err := db.Where("company_id = ? AND employee_id = ?", 1, emp.ID).First(&addr).Error; err != nil {
		t.Fatal("address not found in DB")
	}
	if addr.PostCode == nil || *addr.PostCode != "100-0001" {
		t.Errorf("want 100-0001, got %v", addr.PostCode)
	}
}

func TestUpsertAddress_UpdatesWhenExists(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)
	emp := seedEmployee(t, db, 1, "EA002")

	old := "000-0000"
	db.Create(&model.EmployeeAddress{CompanyId: 1, EmployeeId: emp.ID, PostCode: &old})

	newCode := "100-0001"
	if err := repo.UpsertAddress(&model.EmployeeAddress{
		CompanyId: 1, EmployeeId: emp.ID, PostCode: &newCode,
	}); err != nil {
		t.Fatal(err)
	}

	var count int64
	db.Model(&model.EmployeeAddress{}).Where("company_id = ? AND employee_id = ?", 1, emp.ID).Count(&count)
	if count != 1 {
		t.Errorf("want exactly 1 address record after upsert, got %d", count)
	}

	var addr model.EmployeeAddress
	db.Where("company_id = ? AND employee_id = ?", 1, emp.ID).First(&addr)
	if addr.PostCode == nil || *addr.PostCode != "100-0001" {
		t.Errorf("want updated 100-0001, got %v", addr.PostCode)
	}
}

// --- UpdateDepartments ---

func TestUpdateDepartments_AssignsDepartments(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)

	dept := model.Department{CompanyId: 1, Name: "営業部"}
	db.Create(&dept)

	if err := repo.UpdateDepartments(1, 1, []uint{dept.ID}); err != nil {
		t.Fatal(err)
	}

	var count int64
	db.Model(&model.EmployeeDepartments{}).Where("employee_id = 1").Count(&count)
	if count != 1 {
		t.Errorf("want 1 department assigned, got %d", count)
	}
}

func TestUpdateDepartments_ReplacesPreviousAssignments(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)

	dept1 := model.Department{CompanyId: 1, Name: "営業部"}
	dept2 := model.Department{CompanyId: 1, Name: "開発部"}
	db.Create(&dept1)
	db.Create(&dept2)

	repo.UpdateDepartments(1, 1, []uint{dept1.ID})
	if err := repo.UpdateDepartments(1, 1, []uint{dept2.ID}); err != nil {
		t.Fatal(err)
	}

	var records []model.EmployeeDepartments
	db.Where("employee_id = 1").Find(&records) // GORM auto-filters soft-deleted
	if len(records) != 1 {
		t.Errorf("want 1 active department, got %d", len(records))
	}
	if records[0].DepartmentId != dept2.ID {
		t.Errorf("want dept2 ID=%d, got %d", dept2.ID, records[0].DepartmentId)
	}
}

func TestUpdateDepartments_ClearsWhenEmpty(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)

	dept := model.Department{CompanyId: 1, Name: "営業部"}
	db.Create(&dept)
	repo.UpdateDepartments(1, 1, []uint{dept.ID})

	if err := repo.UpdateDepartments(1, 1, []uint{}); err != nil {
		t.Fatal(err)
	}

	var count int64
	db.Model(&model.EmployeeDepartments{}).Where("employee_id = 1").Count(&count)
	if count != 0 {
		t.Errorf("want 0 departments after clearing, got %d", count)
	}
}

// --- ReplaceTenures ---

func TestReplaceTenures_ReplacesExisting(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)
	emp := seedEmployee(t, db, 1, "RT001")

	oldStatus := "在職中"
	db.Create(&model.EmployeeTenures{
		CompanyId: 1, EmployeeId: emp.ID,
		JoinedOn: parseDate("2018-04-01"), Status: &oldStatus,
	})

	newStatus := "退職済"
	if err := repo.ReplaceTenures(1, emp.ID, []model.EmployeeTenures{
		{CompanyId: 1, EmployeeId: emp.ID, JoinedOn: parseDate("2020-04-01"), Status: &newStatus},
	}); err != nil {
		t.Fatal(err)
	}

	var records []model.EmployeeTenures
	db.Where("employee_id = ?", emp.ID).Find(&records)
	if len(records) != 1 {
		t.Fatalf("want 1 tenure after replace, got %d", len(records))
	}
	if records[0].Status == nil || *records[0].Status != "退職済" {
		t.Errorf("want 退職済, got %v", records[0].Status)
	}
}

func TestReplaceTenures_ClearsWhenEmpty(t *testing.T) {
	db := setupEmployeeDB(t)
	repo := repository.NewEmployeeRepository(db)
	emp := seedEmployee(t, db, 1, "RT002")

	status := "在職中"
	db.Create(&model.EmployeeTenures{
		CompanyId: 1, EmployeeId: emp.ID,
		JoinedOn: parseDate("2020-04-01"), Status: &status,
	})

	if err := repo.ReplaceTenures(1, emp.ID, []model.EmployeeTenures{}); err != nil {
		t.Fatal(err)
	}

	var count int64
	db.Model(&model.EmployeeTenures{}).Where("employee_id = ?", emp.ID).Count(&count)
	if count != 0 {
		t.Errorf("want 0 tenures after clearing, got %d", count)
	}
}
