package model

import (
	"sync"
)

type JobDB struct {
	DAGRunId string `gorm:"column:dagrun_id" json:"dagrun_id"`
	TaskId   string `gorm:"column:task_id" json:"task_id"`
	JobId    string `gorm:"column:job_id; primary_key; auto_increment:true" json:"job_id"`
	Input    string `gorm:"column:intput" json:"input"`
	Output   string `gorm:"column:output" json:"output"`
	Status   uint32 `gorm:"column:status" json:"status"`
}

type JobHTTPReq struct {
	JobID           string            `json:"job_id"`
	IsBoundary      bool              `json:"is_boundary"`
	TaskID          string            `json:"task_id"`
	DAGRunID        string            `json:"dag_run_id"`
	UserID          uint32            `json:"user_id"`
	Command         string            `json:"command"`
	Param           string            `json:"param"`
	ParamsWithRegex []*ParamWithRegex `json:"paramwithregex"`
	ParentJobsID    []string          `json:"parent_jobs_id"`
	ChildrenJobsID  []string          `json:"children_jobs_id"`
	InputLocation   []string          `json:"input_location"`
	DockerImage     []string          `json:"docker_image"`
	QueueLevel      int               `json:"queue_level"`
}

type JobLogic struct {
	JobLock           sync.Mutex
	JobID             string            `gorm:"column:job_id; primaryKey" json:"job_id"`
	IsBoundary        bool              `json:"is_boundary" gorm:"column:is_boundary"`
	TaskID            string            `json:"task_id" gorm:"column:task_of_job"`
	Task              Task              `gorm:"foreignKey:task_of_job;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DAGRunID          string            `json:"dagrun_id" gorm:"column:dagrun_of_job"`
	DAGRun            DAGRun            `gorm:"foreignKey:dagrun_of_job;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID            uint32            `json:"user_id" gorm:"column:user_id"`
	Command           string            `json:"command" gorm:"column:command"`
	ERam              int64             `json:"eram" gorm:"column:eram"`
	ECPU              int64             `json:"ecpu" gorm:"column:ecpu"`
	Param             string            `json:"param" gorm:"column:param"`
	ParamsWithRegex   []*ParamWithRegex `json:"paramwithregex"`
	ParamsWithRegexDB []byte            `json:"paramwithregexdb" gorm:"Column:params_with_regex_db; Type:jsonb"`
	ParentJobs        []*JobLogic       `json:"parent_jobs" gorm:"-"`
	ChildrenJobs      []*JobLogic       `json:"children_jobs" gorm:"-"`
	ParentJobsId      []string          `json:"parents_id" gorm:"type:text[];column:parents_id"`
	ChildrenJobsId    []string          `json:"children_id" gorm:"type:text[];column:children_id"`
	InputLocation     []string          `json:"input_location" gorm:"type:text[];column:input_location"`
	InputFileSize     []int64           `json:"input_file_size" gorm:"column:input_file_size"`
	OutputLocation    []string          `json:"output_location" gorm:"type:text[];column:output_location"`
	DockerImage       []string          `json:"docker_image" gorm:"type:text[];column:docker_image"`
	QueueLevel        int               `json:"queue_level" gorm:"column:queue_level"`
	ParentsDoneCount  uint32            `json:"parent_done_count" gorm:"column:parent_done_count"`
	ProcessStatus     uint32            `json:"process_status" gorm:"column:process_status"` // 0: not processing, 1:ready, 2: processing, 3: done, 4: fail, 5: in kafka queue
}

type ParamWithRegex struct {
	Scatter        bool             `json:"scatter"`
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
