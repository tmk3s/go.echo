package worker

const (
	TypeBulkCreate = "employee:bulk_create"
	TypeBulkUpdate = "employee:bulk_update"
)

type JobPayload struct {
	JobId uint `json:"job_id"`
}
