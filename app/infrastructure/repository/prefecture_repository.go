package repository

import (
	"app/domain/model"
	"app/domain/repository"

	"gorm.io/gorm"
)

type prefectureRepository struct {
	Conn *gorm.DB
}

func NewPrefectureRepository(Conn *gorm.DB) repository.PrefectureRepository {
	return &prefectureRepository{Conn}
}

func (r *prefectureRepository) GetAll() ([]model.Prefecture, error) {
	var prefectures []model.Prefecture
	err := r.Conn.Order("id").Find(&prefectures).Error
	if err != nil {
		return nil, err
	}
	return prefectures, nil
}
