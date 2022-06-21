package model

type SelectTaskReq struct {
	CPUpoint     int64    `json:"cpu_point"`
	RAMpoint     int64    `json:"ram_point"`
	PenaltyList  []uint64 `json:"penalty_list" db:"penalty_list"`
	ExecutorName string   `json:"executor_name"`
}

type ParamWithRegex struct {
	From           []string         `json:"from"`
	SecondaryFiles []string         `json:"secondary_files"`
	Regex          []string         `json:"regex"`
	Files          []*FilteredFiles `json:"files"`
	Prefix         string           `json:"prefix"`
}

type FilteredFiles struct {
	Filepath string `json:"file_path"`
	Filesize int64  `json:"file_size"`
}

type TaskHTTPResp struct {
	TaskID           string            `gorm:"column:task_id; primary_key" json:"task_id"`
	TaskUUID         string            `gorm:"column:task_uuid;" json:"task_uuid"`
	TaskName         string            `json:"task_name"`
	IsBoundary       bool              `json:"is_boundary" gorm:"column:is_boundary"`
	StepID           string            `json:"step_id" gorm:"column:step_id"`
	RunID            int               `json:"run_id" gorm:"column:run_id"`
	RunUUID          string            `json:"run_uuid" gorm:"column:run_uuid"`
	RunName          string            `json:"run_name"`
	UserID           string            `json:"username" gorm:"column:username"`
	Command          string            `json:"command" gorm:"column:command"`
	ERam             int64             `json:"eram" gorm:"column:eram"`
	ECPU             int64             `json:"ecpu" gorm:"column:ecpu"`
	ParamsWithRegex  []*ParamWithRegex `json:"paramwithregex"`
	OutputRegex      []string          `json:"output_regex"`
	Output2ndFiles   []string          `json:"output_2nd_files"`
	DockerImage      []string          `json:"docker_image" gorm:"type:text[];column:docker_image"`
	QueueLevel       int               `json:"queue_level" gorm:"column:queue_level"`
	ParentsDoneCount uint32            `json:"parent_done_count" gorm:"column:parent_done_count"`
	ProcessStatus    uint32            `json:"process_status" gorm:"column:process_status"` // 0: not processing, 1:ready, 2: processing, 3: done, 4: fail
}

type UpdateStatusCheck struct {
	Success  bool     `json:"success"`
	TaskID   string   `json:"task_id"`
	Filename []string `json:"file_name"`
	Filesize []int64  `json:"file_size" gorm:"Type:BIGINT"`
}

type SelectTaskResp struct {
	Ack   string         `json:"ack"`
	Tasks []TaskHTTPResp `json:"tasks"`
}

type DeleteTaskK8SReq struct {
	TaskID []string `json:"task_id"`
}

type DeleteTaskNoti struct {
	TaskID string `json:"task_id"`
}

type ExecutorResumeWithoutTaskReq struct {
	ExecutorName string `json:"executor_name"`
}
