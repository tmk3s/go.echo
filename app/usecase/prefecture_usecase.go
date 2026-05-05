package usecase

import (
	"app/domain/model"
	"app/domain/repository"
)

type PrefectureUseCase interface {
	GetAll() ([]model.Prefecture, error)
}

type prefectureUseCase struct {
	repo repository.PrefectureRepository
}

func NewPrefectureUseCase(r repository.PrefectureRepository) PrefectureUseCase {
	return &prefectureUseCase{repo: r}
}

func (u *prefectureUseCase) GetAll() ([]model.Prefecture, error) {
	return u.repo.GetAll()
}
