package usecase_test

import (
	"testing"

	"app/domain/model"
	"app/usecase"
)

func TestAuthGetUser_ReturnsUserFromRepo(t *testing.T) {
	expected := &model.User{Email: "test@example.com"}
	expected.ID = 1
	repo := &mockUserRepo{user: expected}
	uc := usecase.NewAuthUseCase(repo)

	got, err := uc.GetUser("test@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Email != "test@example.com" {
		t.Errorf("want user with email test@example.com, got %v", got)
	}
}

func TestAuthGetUser_ReturnsNilWhenNotFound(t *testing.T) {
	repo := &mockUserRepo{user: nil}
	uc := usecase.NewAuthUseCase(repo)

	got, err := uc.GetUser("nobody@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Error("want nil for non-existent user, got user")
	}
}
