package dto

type UpdateTaskOutput struct {
	TaskID string   `json:"task_id"`
	Output []string `json:"output_location"`
}

type UpdateTaskStatus struct {
	TaskIDs []string `json:"task_ids"`
	Status  uint32   `json:"status"`
}
