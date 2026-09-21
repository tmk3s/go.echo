package usecase_test

import (
	"testing"

	"app/domain/model"
	"app/usecase"
)

func TestJobGetById_OwnCompany_ReturnsJob(t *testing.T) {
	jobRepo := &mockJobRepo{}
	created, _ := jobRepo.Create(&model.Job{CompanyId: 1, Status: "completed"})

	uc := usecase.NewJobUseCase(jobRepo)
	job, err := uc.GetById(1, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if job.CompanyId != 1 {
		t.Errorf("want company 1, got %d", job.CompanyId)
	}
}

func TestJobGetById_OtherCompany_ReturnsNotFound(t *testing.T) {
	jobRepo := &mockJobRepo{}
	// ErrorMessage に他社の社員コードやメールアドレスが含まれうる
	created, _ := jobRepo.Create(&model.Job{
		CompanyId:    2,
		Status:       "failed",
		ErrorMessage: "S001(taro@example.com) は登録済みです",
	})

	uc := usecase.NewJobUseCase(jobRepo)
	job, err := uc.GetById(1, created.ID)
	if err == nil {
		t.Fatal("want error when fetching another company's job, got nil")
	}
	if job != nil {
		t.Error("another company's job must not be returned")
	}
}
