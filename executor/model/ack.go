package model

type TaskAck struct {
	TaskID string `json:"task_id"`
	Status uint32 `json:"status"`
}
