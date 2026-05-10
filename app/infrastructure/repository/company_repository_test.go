package repository_test

import (
	"testing"

	"gorm.io/gorm"

	"app/domain/model"
	"app/infrastructure/repository"
)

func setupCompanyDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := db.AutoMigrate(&model.Company{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { truncate(t, db, "companies") })
	return db
}

func seedCompany(t *testing.T, db *gorm.DB, name string) model.Company {
	t.Helper()
	c := model.Company{Name: name}
	if err := db.Create(&c).Error; err != nil {
		t.Fatal(err)
	}
	return c
}

// --- GetById ---

func TestCompanyGetById_Found(t *testing.T) {
	db := setupCompanyDB(t)
	repo := repository.NewCompanyRepository(db)
	created := seedCompany(t, db, "株式会社サンプル")

	got, err := repo.GetById(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "株式会社サンプル" {
		t.Errorf("want 株式会社サンプル, got %s", got.Name)
	}
}

func TestCompanyGetById_NotFound(t *testing.T) {
	repo := repository.NewCompanyRepository(setupCompanyDB(t))

	_, err := repo.GetById(9999)
	if err == nil {
		t.Error("want error for non-existent company, got nil")
	}
}

// --- Update ---

func TestCompanyUpdate_ChangesName(t *testing.T) {
	db := setupCompanyDB(t)
	repo := repository.NewCompanyRepository(db)
	created := seedCompany(t, db, "旧社名株式会社")

	created.Name = "新社名株式会社"
	if err := repo.Update(&created); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetById(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "新社名株式会社" {
		t.Errorf("want 新社名株式会社, got %s", got.Name)
	}
}
