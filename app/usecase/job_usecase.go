package usecase

import (
	"app/domain/model"
	"app/domain/repository"
)

type JobUseCase interface {
	GetById(id uint) (*model.Job, error)
}

type jobUseCase struct {
	repo repository.JobRepository
}

func NewJobUseCase(repo repository.JobRepository) JobUseCase {
	return &jobUseCase{repo: repo}
}

func (u *jobUseCase) GetById(id uint) (*model.Job, error) {
	return u.repo.GetById(id)
}
