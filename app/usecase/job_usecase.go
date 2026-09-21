package usecase

import (
	"errors"

	"app/domain/model"
	"app/domain/repository"
)

// ErrJobNotFound はジョブが存在しない、または自社のジョブでない場合に返す。
// 他社のジョブであることを区別して返すと存在自体が漏れるため、同じエラーにまとめている。
var ErrJobNotFound = errors.New("job not found")

type JobUseCase interface {
	GetById(companyId uint, id uint) (*model.Job, error)
}

type jobUseCase struct {
	repo repository.JobRepository
}

func NewJobUseCase(repo repository.JobRepository) JobUseCase {
	return &jobUseCase{repo: repo}
}

func (u *jobUseCase) GetById(companyId uint, id uint) (*model.Job, error) {
	job, err := u.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	// ErrorMessage に他社の社員コードやメールアドレスが含まれるため、会社でスコープする
	if job == nil || job.CompanyId != companyId {
		return nil, ErrJobNotFound
	}
	return job, nil
}
