package worker

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type AsynqEnqueuer struct {
	client *asynq.Client
}

func NewAsynqEnqueuer(client *asynq.Client) *AsynqEnqueuer {
	return &AsynqEnqueuer{client: client}
}

func (e *AsynqEnqueuer) EnqueueBulkCreate(jobId uint) error {
	payload, err := json.Marshal(JobPayload{JobId: jobId})
	if err != nil {
		return err
	}
	task := asynq.NewTask(TypeBulkCreate, payload, asynq.MaxRetry(0))
	_, err = e.client.Enqueue(task)
	return err
}

func (e *AsynqEnqueuer) EnqueueBulkUpdate(jobId uint) error {
	payload, err := json.Marshal(JobPayload{JobId: jobId})
	if err != nil {
		return err
	}
	task := asynq.NewTask(TypeBulkUpdate, payload, asynq.MaxRetry(0))
	_, err = e.client.Enqueue(task)
	return err
}
