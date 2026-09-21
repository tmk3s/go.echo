package repository_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"app/domain/model"
	"app/infrastructure/repository"
)

func setupCompanyDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := db.AutoMigrate(&model.Company{}, &model.User{}, &model.UserInfo{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { truncate(t, db, "user_infos", "users", "companies") })
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

// --- Register ---

func TestCompanyRegister_LinksUserToCompany(t *testing.T) {
	db := setupCompanyDB(t)
	repo := repository.NewCompanyRepository(db)

	company := &model.Company{Name: "株式会社テスト"}
	user := model.NewCompanyOwner("owner@example.com", "password123", "山田", "太郎")

	if err := repo.Register(company, user); err != nil {
		t.Fatal(err)
	}
	if company.ID == 0 {
		t.Fatal("want company to be created with an ID")
	}

	var got model.User
	if err := db.Preload("UserInfo").First(&got, user.ID).Error; err != nil {
		t.Fatal(err)
	}
	// company_id = 0 のユーザーが作られると他社のデータが見えてしまう
	if got.CompanyId != company.ID {
		t.Errorf("want user linked to company %d, got %d", company.ID, got.CompanyId)
	}
	if got.UserInfo.LastName != "山田" || got.UserInfo.FirstName != "太郎" {
		t.Errorf("want 山田 太郎, got %s %s", got.UserInfo.LastName, got.UserInfo.FirstName)
	}
	if got.Password == "password123" {
		t.Error("password must not be stored in plain text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte("password123")); err != nil {
		t.Errorf("stored password does not match the original: %v", err)
	}
}

func TestCompanyRegister_RollsBackWhenUserCreationFails(t *testing.T) {
	db := setupCompanyDB(t)
	repo := repository.NewCompanyRepository(db)

	// 先に同じメールアドレスのユーザーを作っておく(users.email は uniqueIndex)
	existingCompany := seedCompany(t, db, "既存の会社")
	if err := db.Create(&model.User{CompanyId: existingCompany.ID, Email: "dup@example.com", Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}

	company := &model.Company{Name: "登録されないはずの会社"}
	user := model.NewCompanyOwner("dup@example.com", "password123", "山田", "太郎")
	if err := repo.Register(company, user); err == nil {
		t.Fatal("want error for duplicate email, got nil")
	}

	var count int64
	db.Model(&model.Company{}).Where("name = ?", "登録されないはずの会社").Count(&count)
	if count != 0 {
		t.Error("company must be rolled back when the user cannot be created")
	}
}
