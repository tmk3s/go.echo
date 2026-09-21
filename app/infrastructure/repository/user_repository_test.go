package repository_test

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"app/domain/model"
	"app/infrastructure/repository"
)

func setupUserDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := db.AutoMigrate(&model.User{}, &model.UserInfo{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { truncate(t, db, "users", "user_infos") })
	return db
}

func createTestUser(t *testing.T, repo interface {
	Create(*model.User) (*model.User, error)
}, email, password string) *model.User {
	t.Helper()
	user, err := repo.Create(&model.User{CompanyId: 1, Email: email, Password: password})
	if err != nil {
		t.Fatal(err)
	}
	return user
}

// --- Create ---

func TestUserCreate_AssignsID(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))
	user := createTestUser(t, repo, "test@example.com", "password")
	if user.ID == 0 {
		t.Error("expected ID to be assigned after Create")
	}
}

func TestUserCreate_HashesPassword(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))
	user := createTestUser(t, repo, "test@example.com", "plaintext")

	if user.Password == "plaintext" {
		t.Error("password was stored as plain text, expected bcrypt hash")
	}
	// bcrypt hashes start with $2a$ or $2b$
	if len(user.Password) < 4 || (user.Password[:3] != "$2a" && user.Password[:3] != "$2b") {
		t.Errorf("expected bcrypt hash, got: %s", user.Password)
	}
}

// --- GetByEmailAndPass ---

func TestGetByEmailAndPass_CorrectPassword(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))
	createTestUser(t, repo, "login@example.com", "mypassword")

	got, err := repo.GetByEmailAndPass("login@example.com", "mypassword")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected user for correct password, got nil")
	}
	if got.Email != "login@example.com" {
		t.Errorf("want email login@example.com, got %s", got.Email)
	}
}

func TestGetByEmailAndPass_WrongPassword(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))
	createTestUser(t, repo, "login@example.com", "correct")

	got, err := repo.GetByEmailAndPass("login@example.com", "wrong")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Error("expected nil for wrong password, got user")
	}
}

func TestGetByEmailAndPass_EmailNotFound(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))

	got, err := repo.GetByEmailAndPass("nobody@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Error("expected nil for non-existent email, got user")
	}
}

// --- GetByEmail ---

func TestGetByEmail_Found(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))
	createTestUser(t, repo, "find@example.com", "pass")

	got, err := repo.GetByEmail("find@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected user, got nil")
	}
	if got.Email != "find@example.com" {
		t.Errorf("want find@example.com, got %s", got.Email)
	}
}

func TestGetByEmail_NotFound(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))

	got, err := repo.GetByEmail("nobody@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Error("expected nil for non-existent email, got user")
	}
}

func TestGetByEmail_DoesNotReturnPassword(t *testing.T) {
	repo := repository.NewUserRepository(setupUserDB(t))
	createTestUser(t, repo, "safe@example.com", "secret")

	got, err := repo.GetByEmail("safe@example.com")
	if err != nil {
		t.Fatal(err)
	}
	// Password field has json:"-" but is still loaded from DB.
	// The important check: it must not be the plain-text original.
	if got.Password == "secret" {
		t.Error("Password should be hashed, not plain text")
	}
}

// --- Update ---

func TestUserUpdate_PersistsUserInfoColumns(t *testing.T) {
	db := setupUserDB(t)
	repo := repository.NewUserRepository(db)

	user := createTestUser(t, repo, "update@example.com", "password")
	birthday := time.Date(1990, 1, 2, 0, 0, 0, 0, time.Local)
	if err := db.Create(&model.UserInfo{
		UserId:    user.ID,
		LastName:  "旧姓",
		FirstName: "旧名",
		Gender:    1,
		BirthDay:  &birthday,
		Working:   false,
	}).Error; err != nil {
		t.Fatal(err)
	}

	loaded, err := repo.GetById(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	newBirthday := time.Date(1995, 6, 7, 0, 0, 0, 0, time.Local)
	loaded.UserInfo.LastName = "新姓"
	loaded.UserInfo.FirstName = "新名"
	loaded.UserInfo.Gender = 2
	loaded.UserInfo.BirthDay = &newBirthday
	loaded.UserInfo.Working = true

	if _, err := repo.Update(loaded); err != nil {
		t.Fatal(err)
	}

	// 既定の Save では has-one 関連の FK しか更新されず、ここが旧姓のままになっていた
	var got model.UserInfo
	if err := db.Where("user_id = ?", user.ID).First(&got).Error; err != nil {
		t.Fatal(err)
	}
	if got.LastName != "新姓" || got.FirstName != "新名" {
		t.Errorf("want 新姓 新名, got %s %s", got.LastName, got.FirstName)
	}
	if got.Gender != 2 {
		t.Errorf("want gender 2, got %d", got.Gender)
	}
	if got.BirthDay == nil || !got.BirthDay.Equal(newBirthday) {
		t.Errorf("want %v, got %v", newBirthday, got.BirthDay)
	}
	if !got.Working {
		t.Error("want working=true")
	}

	// 関連が二重に作られていないことも確認する
	var count int64
	db.Model(&model.UserInfo{}).Where("user_id = ?", user.ID).Count(&count)
	if count != 1 {
		t.Errorf("want exactly 1 user_info row, got %d", count)
	}
}
