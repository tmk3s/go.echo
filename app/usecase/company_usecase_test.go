package usecase_test

import (
	"testing"

	"app/domain/model"
	"app/usecase"
)

// ---- GetCompany ----

func TestGetCompany_ReturnsCompany(t *testing.T) {
	company := &model.Company{Name: "株式会社サンプル"}
	company.ID = 1
	repo := &mockCompanyRepo{company: company}
	uc := usecase.NewCompanyUseCase(repo)

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
	uc := usecase.NewCompanyUseCase(repo)

	if err := uc.UpdateCompany(1, "新社名"); err != nil {
		t.Fatal(err)
	}
	if repo.updatedName != "新社名" {
		t.Errorf("want updated name 新社名, got %s", repo.updatedName)
	}
}
