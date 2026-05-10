package repository

import "app/domain/model"

type JobRepository interface {
	Create(job *model.Job) (*model.Job, error)
	GetById(id uint) (*model.Job, error)
	StartProcessing(id uint, totalCount int) error
	UpdateProgress(id uint, processedCount int) error
	Complete(id uint) error
	Fail(id uint, errMsg string) error
}
