package model

import "time"

type GetMutualFileResponse struct {
	Total int          `json:"total"`
	File  []MutualFile `json:"files"`
}

type GetGeneratedFileResponse struct {
	Total int64           `json:"total"`
	File  []GeneratedFile `json:"files"`
}

type GetUploadedFileResponse struct {
	Total int64          `json:"total"`
	File  []UploadedFile `json:"files"`
}

type GetSampleResponse struct {
	Total   int64               `json:"total"`
	Samples []*SampleDetailResp `json:"samples"`
}

type SampleDetailResp struct {
	SampleUUID   string       `json:"sample_uuid" gorm:"Column:sample_uuid; Type:varchar; Size:255; primaryKey"`
	DatasetUUID  string       `json:"dataset_uuid" gorm:"Column:dataset_uuid; Type:varchar; Size:255"`
	UserID       string       `json:"username" gorm:"Column:username; Type:varchar; Size:255"`
	SampleName   string       `json:"sample_name" gorm:"Column:sample_name; Type:varchar; Size:255"`
	WorkflowUUID string       `json:"workflow_uuid" gorm:"Column:workflow_uuid; Type:varchar; Size:255"`
	SampleFiles  []MutualFile `json:"sample_files"`

	CreatedAt time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`
}

type MutualFile struct {
	FileUUID string `json:"file_uuid" gorm:"Column:file_uuid; Type:varchar; Size:255"`
	// user metadata
	UserID      string `json:"username" gorm:"Column:username; Type:varchar; Size:255"`
	SampleIndex int    `json:"sample_index" gorm:"Column:sample_index; Type:varchar; Size:255"`

	Bucket string `json:"bucket" gorm:"Column:bucket; Type:varchar; Size:255"`
	// local path on hard disk or s3 path
	Path     string `json:"path" gorm:"Column:path; Type:varchar; Size:255"`
	Filename string `json:"file_name" gorm:"Column:filename; Type:varchar; Size:255"`
	Filesize int64  `json:"file_size" gorm:"Column:filesize; Type:bigint; Size:255"`

	CreatedAt time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`
}

type CreateNewSampleReq struct {
	SampleName   string   `json:"sample_name"`
	DatasetUUID  string   `json:"dataset_uuid"`
	WorkflowUUID string   `json:"workflow_uuid"`
	FileUUIDs    []string `json:"file_uuid"`
}

type CreateNewDatasetReq struct {
	DatasetName string `json:"dataset_name"`
}

type CreatePresignedURLReq struct {
	Username string `json:"username"`
	RunUUID  string `json:"run_uuid"`
	TaskUUID string `json:"task_uuid"`
	Filename string `json:"file_name"`
	FileUUID string `json:"file_uuid"`
	TTL      int    `json:"ttl"`
}

type DownloadFileReq struct {
	RunUUID  string `json:"run_uuid"`
	TaskUUID string `json:"task_uuid"`
	Filename string `json:"file_name"`
}

type DeleteFailRunDirReq struct {
	RunUUID  string   `json:"run_uuid"` // if this field != "" then this is a fail run notification request
	UserID   string   `json:"username"`
	TaskID   string   `json:"task_id"`
	Filename []string `json:"file_name"`
}

type HandlerGeneratedFileReq struct {
	UserID      string   `json:"username"`
	RunID       int      `json:"run_id"`
	RunUUID     string   `json:"run_uuid"`
	RunName     string   `json:"run_name"`
	ProjectID   string   `json:"project_id" gorm:"Column:project_id; Type:varchar; Size:255"`
	ProjectPath string   `json:"project_path" gorm:"Column:project_path;"`
	Taskid      string   `json:"task_id"`
	TaskUUID    string   `json:"task_uuid"`
	TaskName    string   `json:"task_name"`
	Filename    []string `json:"file_name"`
	Filesize    []int64  `json:"file_size"`
	LastTask    bool     `json:"last_task"`
}

type FileUploadResp struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type ReuploadToMinioReq struct {
	UserID   string `json:"username"`
	FileUUID string `json:"file_uuid"`
	RunID    string `json:"run_id"`
	RunUUID  string `json:"run_uuid"`
	RunName  string `json:"run_name"`
	Taskid   string `json:"task_id"`
	TaskUUID string `json:"task_uuid"`
	TaskName string `json:"task_name"`
	Filename string `json:"file_name"`
	LastTask bool   `json:"last_task"`
}

type UpdateFilePath struct {
	FileID      string `json:"file_id"`
	ProjectPath string `json:"project_path"`
}

type UpdatePathFiles struct {
	PathFiles []UpdateFilePath `json:"path_files"`
}
