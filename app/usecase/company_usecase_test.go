package usecase_test

import (
	"strings"
	"testing"

	"app/domain/model"
	"app/usecase"
)

// ---- GetCompany ----

func TestGetCompany_ReturnsCompany(t *testing.T) {
	company := &model.Company{Name: "株式会社サンプル"}
	company.ID = 1
	repo := &mockCompanyRepo{company: company}
	uc := usecase.NewCompanyUseCase(repo, &mockUserRepo{})

	got, err := uc.GetCompany(1)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Name != "株式会社サンプル" {
		t.Errorf("want 株式会社サンプル, got %v", got)
	}
}

// ---- UpdateCompany ----

func TestUpdateCompany_ChangesName(t *testing.T) {
	company := &model.Company{Name: "旧社名"}
	company.ID = 1
	repo := &mockCompanyRepo{company: company}
	uc := usecase.NewCompanyUseCase(repo, &mockUserRepo{})

	if err := uc.UpdateCompany(1, "新社名"); err != nil {
		t.Fatal(err)
	}
	if repo.updatedName != "新社名" {
		t.Errorf("want updated name 新社名, got %s", repo.updatedName)
	}
}

// ---- Register ----

func validRegistrationInput() usecase.CompanyRegistrationInput {
	return usecase.CompanyRegistrationInput{
		CompanyName: "株式会社テスト",
		Email:       "owner@example.com",
		Password:    "password123",
		LastName:    "山田",
		FirstName:   "太郎",
	}
}

func TestCompanyRegister_CreatesCompanyAndUser(t *testing.T) {
	companyRepo := &mockCompanyRepo{}
	userRepo := &mockUserRepo{user: nil} // 未登録のメールアドレス
	uc := usecase.NewCompanyUseCase(companyRepo, userRepo)

	company, err := uc.Register(validRegistrationInput())
	if err != nil {
		t.Fatal(err)
	}
	if company == nil || company.ID == 0 {
		t.Fatal("want created company with ID")
	}
	if company.Name != "株式会社テスト" {
		t.Errorf("want 株式会社テスト, got %s", company.Name)
	}

	user := companyRepo.registeredUser
	if user == nil {
		t.Fatal("want the first user to be created together with the company")
	}
	// company_id = 0 のユーザーが作られると他社のデータが見えてしまうため、
	// 会社に紐づいていることを必ず確認する
	if user.CompanyId != company.ID {
		t.Errorf("want user linked to company %d, got %d", company.ID, user.CompanyId)
	}
	if user.Email != "owner@example.com" {
		t.Errorf("want owner@example.com, got %s", user.Email)
	}
	if user.UserInfo.LastName != "山田" || user.UserInfo.FirstName != "太郎" {
		t.Errorf("want 山田 太郎, got %s %s", user.UserInfo.LastName, user.UserInfo.FirstName)
	}
}

func TestCompanyRegister_DuplicateEmail(t *testing.T) {
	existing := &model.User{Email: "owner@example.com"}
	existing.ID = 1
	companyRepo := &mockCompanyRepo{}
	uc := usecase.NewCompanyUseCase(companyRepo, &mockUserRepo{user: existing})

	if _, err := uc.Register(validRegistrationInput()); err == nil {
		t.Fatal("want error for duplicate email, got nil")
	}
	if companyRepo.registeredUser != nil {
		t.Error("company must not be created when the email is already taken")
	}
}

func TestCompanyRegister_ValidatesInput(t *testing.T) {
	cases := map[string]func(*usecase.CompanyRegistrationInput){
		"会社名が空":     func(i *usecase.CompanyRegistrationInput) { i.CompanyName = "  " },
		"メールアドレスが空": func(i *usecase.CompanyRegistrationInput) { i.Email = "" },
		"パスワードが短い":  func(i *usecase.CompanyRegistrationInput) { i.Password = "short" },
		"姓が空":       func(i *usecase.CompanyRegistrationInput) { i.LastName = "" },
		"名が空":       func(i *usecase.CompanyRegistrationInput) { i.FirstName = "" },
	}

	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			input := validRegistrationInput()
			mutate(&input)

			companyRepo := &mockCompanyRepo{}
			uc := usecase.NewCompanyUseCase(companyRepo, &mockUserRepo{})
			if _, err := uc.Register(input); err == nil {
				t.Fatal("want validation error, got nil")
			}
			if companyRepo.registeredUser != nil {
				t.Error("company must not be created when the input is invalid")
			}
		})
	}
}

func TestCompanyRegister_TrimsWhitespace(t *testing.T) {
	companyRepo := &mockCompanyRepo{}
	uc := usecase.NewCompanyUseCase(companyRepo, &mockUserRepo{})

	input := validRegistrationInput()
	input.CompanyName = "  株式会社テスト  "
	input.Email = "  owner@example.com  "

	company, err := uc.Register(input)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(company.Name) != company.Name {
		t.Errorf("company name should be trimmed, got %q", company.Name)
	}
	if companyRepo.registeredUser.Email != "owner@example.com" {
		t.Errorf("email should be trimmed, got %q", companyRepo.registeredUser.Email)
	}
}
