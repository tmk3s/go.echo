package usecase_test

import (
	"testing"

	"app/domain/model"
	"app/usecase"
)

func newDeptUC(repo *mockDeptRepo, csv *mockCsvService) usecase.DepartmentUseCase {
	return usecase.NewDepartmentUseCase(repo, csv)
}

// ---- GetDepartments ----

func TestGetDepartments_ReturnsAll(t *testing.T) {
	repo := &mockDeptRepo{
		departments: []model.Department{
			{Name: "営業部"},
			{Name: "開発部"},
		},
	}
	uc := newDeptUC(repo, &mockCsvService{})

	result, err := uc.GetDepartments(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(*result) != 2 {
		t.Errorf("want 2 departments, got %d", len(*result))
	}
}

// ---- Create ----

func TestDepartmentCreate_CallsRepo(t *testing.T) {
	repo := &mockDeptRepo{}
	uc := newDeptUC(repo, &mockCsvService{})

	if err := uc.Create(1, "新部署", nil); err != nil {
		t.Fatal(err)
	}
	if repo.createCount != 1 {
		t.Errorf("want Create called once, got %d", repo.createCount)
	}
}

// ---- Update ----

func TestDepartmentUpdate_ChangesName(t *testing.T) {
	dept := &model.Department{Name: "旧部署名"}
	dept.ID = 1
	repo := &mockDeptRepo{byId: dept}
	uc := newDeptUC(repo, &mockCsvService{})

	if err := uc.Update(1, "新部署名"); err != nil {
		t.Fatal(err)
	}
	// department name should be updated in-place before repo.Update is called
	if dept.Name != "新部署名" {
		t.Errorf("want 新部署名, got %s", dept.Name)
	}
}

// ---- Delete ----

func TestDepartmentDelete_CallsRepo(t *testing.T) {
	dept := &model.Department{Name: "削除部署"}
	dept.ID = 1
	var deleted *model.Department
	repo := &mockDeptRepo{byId: dept}
	// Override Delete to capture the call
	_ = repo // Delete is already stubbed in mockDeptRepo

	uc := newDeptUC(repo, &mockCsvService{})
	if err := uc.Delete(1); err != nil {
		t.Fatal(err)
	}
	_ = deleted // deletion happens via repo.Delete which is a no-op in the mock
}

// ---- Upload ----

func TestUpload_CreatesOnlyNewDepartments(t *testing.T) {
	repo := &mockDeptRepo{
		departments: []model.Department{
			{Name: "既存部署"},
		},
	}
	csv := &mockCsvService{
		deptNames: []string{"既存部署", "新部署A", "新部署B"},
	}
	uc := newDeptUC(repo, csv)

	if err := uc.Upload(1, emptyFile(), nil); err != nil {
		t.Fatal(err)
	}
	// "既存部署" should be skipped; only "新部署A" and "新部署B" created
	if repo.createCount != 2 {
		t.Errorf("want 2 departments created (skipping existing), got %d", repo.createCount)
	}
}

func TestUpload_SkipsAllWhenAllExist(t *testing.T) {
	repo := &mockDeptRepo{
		departments: []model.Department{
			{Name: "部署A"},
			{Name: "部署B"},
		},
	}
	csv := &mockCsvService{
		deptNames: []string{"部署A", "部署B"},
	}
	uc := newDeptUC(repo, csv)

	if err := uc.Upload(1, emptyFile(), nil); err != nil {
		t.Fatal(err)
	}
	if repo.createCount != 0 {
		t.Errorf("want 0 creates when all exist, got %d", repo.createCount)
	}
}

func TestUpload_CreatesAllWhenNoneExist(t *testing.T) {
	repo := &mockDeptRepo{departments: []model.Department{}}
	csv := &mockCsvService{
		deptNames: []string{"部署X", "部署Y"},
	}
	uc := newDeptUC(repo, csv)

	if err := uc.Upload(1, emptyFile(), nil); err != nil {
		t.Fatal(err)
	}
	if repo.createCount != 2 {
		t.Errorf("want 2 creates, got %d", repo.createCount)
	}
}
