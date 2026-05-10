package repository_test

import (
	"testing"

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
