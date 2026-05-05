package repository

import "app/domain/model"

type PrefectureRepository interface {
	GetAll() ([]model.Prefecture, error)
}
