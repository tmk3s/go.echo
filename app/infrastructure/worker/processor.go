package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"app/usecase"
)

type employeeProcessor struct {
	employeeUC usecase.EmployeeUseCase
}

func newEmployeeProcessor(uc usecase.EmployeeUseCase) *employeeProcessor {
	return &employeeProcessor{employeeUC: uc}
}

func (p *employeeProcessor) handleBulkCreate(ctx context.Context, task *asynq.Task) error {
	var payload JobPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	return p.employeeUC.ProcessBulkCreate(ctx, payload.JobId)
}

func (p *employeeProcessor) handleBulkUpdate(ctx context.Context, task *asynq.Task) error {
	var payload JobPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	return p.employeeUC.ProcessBulkUpdate(ctx, payload.JobId)
}
