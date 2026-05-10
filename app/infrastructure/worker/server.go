package worker

import (
	"github.com/hibiken/asynq"

	"app/usecase"
)

func NewServer(redisAddr string) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 5},
	)
}

func NewServeMux(employeeUC usecase.EmployeeUseCase) *asynq.ServeMux {
	proc := newEmployeeProcessor(employeeUC)
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeBulkCreate, proc.handleBulkCreate)
	mux.HandleFunc(TypeBulkUpdate, proc.handleBulkUpdate)
	return mux
}
