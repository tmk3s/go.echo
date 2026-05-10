package repository_test

import (
	"testing"

	"gorm.io/gorm"

	"app/domain/model"
	"app/infrastructure/repository"
)

func setupDeptDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := db.AutoMigrate(&model.Department{}, &model.DepartmentPath{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { truncate(t, db, "department_paths", "departments") })
	return db
}

// --- GetById ---

func TestDepartmentGetById_Found(t *testing.T) {
	db := setupDeptDB(t)
	repo := repository.NewDepartmentRepository(db)

	dept := model.Department{CompanyId: 1, Name: "営業部", OrderNo: 1}
	db.Create(&dept)

	got, err := repo.GetById(dept.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "営業部" {
		t.Errorf("want 営業部, got %s", got.Name)
	}
}

func TestDepartmentGetById_NotFound(t *testing.T) {
	repo := repository.NewDepartmentRepository(setupDeptDB(t))

	_, err := repo.GetById(9999)
	if err == nil {
		t.Error("want error for non-existent department, got nil")
	}
}

// --- GetList ---

func TestDepartmentGetList_FiltersByCompany(t *testing.T) {
	db := setupDeptDB(t)
	repo := repository.NewDepartmentRepository(db)

	db.Create(&model.Department{CompanyId: 1, Name: "営業部", OrderNo: 1})
	db.Create(&model.Department{CompanyId: 1, Name: "開発部", OrderNo: 2})
	db.Create(&model.Department{CompanyId: 2, Name: "他社部署", OrderNo: 1})

	list, err := repo.GetList(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("want 2 departments for company 1, got %d", len(list))
	}
	// GetList uses SELECT id,name,depth so CompanyId is not populated;
	// the count of 2 (vs 3 total) proves the WHERE company_id filter worked.
	for _, d := range list {
		if d.Name == "他社部署" {
			t.Error("got department from a different company in results")
		}
	}
}

func TestDepartmentGetList_OrderedByOrderNo(t *testing.T) {
	db := setupDeptDB(t)
	repo := repository.NewDepartmentRepository(db)

	// Insert in reverse order
	db.Create(&model.Department{CompanyId: 1, Name: "後", OrderNo: 2})
	db.Create(&model.Department{CompanyId: 1, Name: "前", OrderNo: 1})

	list, err := repo.GetList(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) < 2 {
		t.Fatal("expected at least 2 departments")
	}
	if list[0].Name != "前" {
		t.Errorf("want first=前 (order_no=1), got %s", list[0].Name)
	}
	if list[1].Name != "後" {
		t.Errorf("want second=後 (order_no=2), got %s", list[1].Name)
	}
}

// --- Create (root, no parent) ---

func TestDepartmentCreate_RootDepartment(t *testing.T) {
	db := setupDeptDB(t)
	repo := repository.NewDepartmentRepository(db)

	dept, err := repo.Create(model.NewDepartment(1, "新部署"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if dept.ID == 0 {
		t.Error("expected ID to be assigned after Create")
	}

	// verify DepartmentPath self-reference was created
	var count int64
	db.Model(&model.DepartmentPath{}).
		Where("ancestor_id = ? AND descendant_id = ?", dept.ID, dept.ID).
		Count(&count)
	if count != 1 {
		t.Errorf("want 1 self-referencing DepartmentPath, got %d", count)
	}
}

// --- Update ---

func TestDepartmentUpdate_ChangesName(t *testing.T) {
	db := setupDeptDB(t)
	repo := repository.NewDepartmentRepository(db)

	dept := model.Department{CompanyId: 1, Name: "旧部署名", OrderNo: 1}
	db.Create(&dept)

	dept.Name = "新部署名"
	if _, err := repo.Update(&dept); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetById(dept.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "新部署名" {
		t.Errorf("want 新部署名, got %s", got.Name)
	}
}
