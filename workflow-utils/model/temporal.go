package model

import "time"

const (
	//Workflow names
	DoneTasktWfName   = "DoneTaskWf"
	FailTasktWfName   = "FailTasktWf"
	ExecuteTaskWfName = "ExecuteTaskWf"

	//Activity names
	ExecuteTaskActName = "ExecuteTaskAct"
	DeleteTaskActName  = "DeleteTaskAct"

	SaveGeneratedFileActName = "SaveGeneratedFileAct"

	UpdateTaskSuccessActName   = "UpdateTaskSuccessAct"
	UpdateTaskFailActName      = "UpdateTaskFailAct"
	UpdateTaskStartTimeActName = "UpdateTaskStartTimeAct"
	UpdateTaskStatusActName    = "UpdateTaskStatusAct"

	//QueuName
	ExecutorInternalQueueName = "ExecutorInternalQueue"
)

type ExecuteTaskParam struct {
	Task TaskDTO
}

type ExecuteTaskResult struct {
	TimeStamp time.Time
	Created   bool
	// maybe add executor name
}

type UpdateTaskSuccessParam struct {
	TaskID   string
	Filename []string
	Filesize []int64
}

type UpdateTaskSuccessResult struct {
	UserName    string `json:"username"`
	RunID       int    `json:"run_id"`
	RunUUID     string `json:"run_uuid"`
	RunName     string `json:"run_name"`
	ProjectID   string `json:"project_id" gorm:"Column:project_id; Type:varchar; Size:255"`
	ProjectPath string `json:"project_path" gorm:"Column:project_path;"`
	TaskID      string `json:"task_id"`
	TaskUUID    string `json:"task_uuid"`
	TaskName    string `json:"task_name"`
	Path        []string
	Filename    []string `json:"file_name"`
	Filesize    []int64  `json:"file_size"`
	LastTask    bool     `json:"last_task"`
}

type UpdateTaskFailParam struct {
	TaskID string
}

type UpdateTaskFailResult struct {
}
