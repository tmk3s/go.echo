package repository

import (
	"app/domain/model"
	"app/domain/repository"
	"gorm.io/gorm"
)

type jobRepository struct {
	Conn *gorm.DB
}

func NewJobRepository(conn *gorm.DB) repository.JobRepository {
	return &jobRepository{Conn: conn}
}

func (r *jobRepository) Create(job *model.Job) (*model.Job, error) {
	if err := r.Conn.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (r *jobRepository) GetById(id uint) (*model.Job, error) {
	var job model.Job
	if err := r.Conn.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) StartProcessing(id uint, totalCount int) error {
	return r.Conn.Model(&model.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "processing",
		"total_count": totalCount,
	}).Error
}

func (r *jobRepository) UpdateProgress(id uint, processedCount int) error {
	return r.Conn.Model(&model.Job{}).Where("id = ?", id).Update("processed_count", processedCount).Error
}

func (r *jobRepository) Complete(id uint) error {
	return r.Conn.Model(&model.Job{}).Where("id = ?", id).Update("status", "completed").Error
}

func (r *jobRepository) Fail(id uint, errMsg string) error {
	return r.Conn.Model(&model.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        "failed",
		"error_message": errMsg,
	}).Error
}
