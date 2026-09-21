package usecase_test

import (
	"testing"
	"time"

	"app/domain/model"
	"app/usecase"
)

func TestUserUpdate_PersistsUserInfo(t *testing.T) {
	user := &model.User{Email: "old@example.com"}
	user.ID = 1
	userRepo := &mockUserRepo{user: user}
	uc := usecase.NewUserUseCase(userRepo)

	birthday := time.Date(1990, 1, 2, 0, 0, 0, 0, time.UTC)
	params := usecase.UserParams{
		LastName:  ptr("山田"),
		FirstName: ptr("太郎"),
		Email:     ptr("new@example.com"),
		Gender:    ptr(1),
		BirthDay:  &birthday,
		Working:   ptr(true),
	}
	if err := uc.Update(1, params); err != nil {
		t.Fatal(err)
	}

	// 以前は Update が何も書き込まずに 200 を返していた
	if userRepo.updatedUser == nil {
		t.Fatal("Update was not called on user repository")
	}
	got := userRepo.updatedUser
	if got.Email != "new@example.com" {
		t.Errorf("want new@example.com, got %s", got.Email)
	}
	if got.UserInfo.LastName != "山田" || got.UserInfo.FirstName != "太郎" {
		t.Errorf("want 山田 太郎, got %s %s", got.UserInfo.LastName, got.UserInfo.FirstName)
	}
	if got.UserInfo.UserId != 1 {
		t.Errorf("want UserId=1, got %d", got.UserInfo.UserId)
	}
	if !got.UserInfo.Working {
		t.Error("want Working=true")
	}
	if got.UserInfo.BirthDay == nil || !got.UserInfo.BirthDay.Equal(birthday) {
		t.Errorf("want %v, got %v", birthday, got.UserInfo.BirthDay)
	}
}

func TestUserUpdate_UserNotFound_ReturnsError(t *testing.T) {
	userRepo := &mockUserRepo{user: nil}
	uc := usecase.NewUserUseCase(userRepo)

	if err := uc.Update(999, usecase.UserParams{}); err == nil {
		t.Fatal("want error when user not found, got nil")
	}
}

func ptr[T any](v T) *T { return &v }

func TestUserUpdate_OmittedFieldsAreKept(t *testing.T) {
	existingBirthday := time.Date(1980, 5, 6, 0, 0, 0, 0, time.UTC)
	user := &model.User{
		Email: "keep@example.com",
		UserInfo: model.UserInfo{
			LastName:  "既存",
			FirstName: "花子",
			Gender:    2,
			BirthDay:  &existingBirthday,
			Working:   true,
		},
	}
	user.ID = 1
	userRepo := &mockUserRepo{user: user}
	uc := usecase.NewUserUseCase(userRepo)

	// メールアドレスだけを送る。他の項目が空値で潰されてはいけない。
	if err := uc.Update(1, usecase.UserParams{Email: ptr("new@example.com")}); err != nil {
		t.Fatal(err)
	}

	got := userRepo.updatedUser
	if got == nil {
		t.Fatal("Update was not called on user repository")
	}
	if got.Email != "new@example.com" {
		t.Errorf("want new@example.com, got %s", got.Email)
	}
	if got.UserInfo.LastName != "既存" || got.UserInfo.FirstName != "花子" {
		t.Errorf("name must be kept, got %s %s", got.UserInfo.LastName, got.UserInfo.FirstName)
	}
	if got.UserInfo.Gender != 2 {
		t.Errorf("gender must be kept, got %d", got.UserInfo.Gender)
	}
	if got.UserInfo.BirthDay == nil || !got.UserInfo.BirthDay.Equal(existingBirthday) {
		t.Errorf("birthday must be kept, got %v", got.UserInfo.BirthDay)
	}
	if !got.UserInfo.Working {
		t.Error("working must be kept as true")
	}
}
