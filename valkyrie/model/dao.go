package model

import (
	"time"
)

// TODO: not used, remove this after 06/2021 release
type File struct {
	LocalPath       string    `json:"local_path" gorm:"Column:local_path; Type:varchar; Size:255"`
	SampleName      string    `json:"sample_name" gorm:"Column:sample_name; Type:varchar; Size:255"`
	UserUploadName  string    `json:"user_upload_name" gorm:"Column:user_upload_name; Type:varchar; Size:255"`
	UserID          string    `json:"username" gorm:"Column:username; Type:varchar; Size:255"`
	RunID           int       `json:"run_id" gorm:"Column:run_id; Type:int; Size:255"`
	RunUUID         string    `json:"run_uuid" gorm:"Column:run_uuid; Type:varchar; Size:255"`
	RunName         string    `json:"run_name" gorm:"Column:run_name; Type:varchar; Size:255"`
	TaskID          string    `json:"task_id" gorm:"Column:task_id; Type:varchar; Size:255"`
	TaskUUID        string    `json:"task_uuid" gorm:"Column:task_uuid; Type:varchar; Size:255"`
	TaskName        string    `json:"task_name" gorm:"Column:task_name; Type:varchar; Size:255"`
	Bucket          string    `json:"bucket" gorm:"Column:bucket; Type:varchar; Size:255"`
	Filename        string    `json:"name" gorm:"Column:filename; Type:varchar; Size:255"`
	Filesize        int64     `json:"size" gorm:"Column:filesize; Type:bigint; Size:255"`
	Deleted         bool      `json:"deleted" gorm:"Column:deleted; Type:boolean"`
	Safe            bool      `json:"safe" gorm:"Column:safe; Type:boolean"`
	UploadSuccess   bool      `json:"upload_success" gorm:"Column:upload_success; Type:boolean"`
	DoneRun         bool      `json:"done_run" gorm:"Column:done_run; Type:boolean"`
	WorkflowUUID    string    `json:"workflow_uuid" gorm:"workflow_uuid"`
	UploadExpiredAt time.Time `json:"upload_expired_at" gorm:"Column:upload_expired_at; Type:timestamp"`
	CreatedAt       time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`
	ExpiredAt       time.Time `json:"expired_at" gorm:"Column:expired_at; Type:timestamp"`
}

type Dataset struct {
	DatasetUUID string `json:"dataset_uuid" gorm:"Column:dataset_uuid; Type:varchar; Size:255; primaryKey"`
	UserID      string `json:"username" gorm:"Column:username; Type:varchar; Size:255"`
	DatasetName string `json:"dataset_name" gorm:"Column:dataset_name; Type:varchar; Size:255"`

	CreatedAt time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`
}

type Sample struct {
	SampleUUID   string `json:"sample_uuid" gorm:"Column:sample_uuid; Type:varchar; Size:255; primaryKey"`
	DatasetUUID  string `json:"dataset_uuid" gorm:"Column:dataset_uuid; Type:varchar; Size:255"`
	UserID       string `json:"username" gorm:"Column:username; Type:varchar; Size:255"`
	SampleName   string `json:"sample_name" gorm:"Column:sample_name; Type:varchar; Size:255"`
	WorkflowUUID string `json:"workflow_uuid" gorm:"Column:workflow_uuid; Type:varchar; Size:255"`

	CreatedAt time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`
}

type SampleContent struct {
	FileUUID    string `json:"file_uuid" gorm:"Column:file_uuid; Type:varchar; Size:255"`
	SampleUUID  string `json:"sample_uuid" gorm:"Column:sample_uuid; Type:varchar; Size:255"`
	SampleIndex int    `json:"sample_index" gorm:"Column:sample_index; Size:255"`

	//UploadFileFK   UploadedFile  `gorm:"foreignKey:FileUUID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	//GenerateFileFK GeneratedFile `gorm:"foreignKey:FileUUID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SampleFK Sample `gorm:"foreignKey:SampleUUID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DataFile struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	ProjectPath string    `json:"project_path"`
	Path        string    `json:"path"`
	Filename    string    `json:"name"`
	Filesize    int64     `json:"size"`
	Owner       string    `json:"owner"`
	CreatedAt   time.Time `json:"created_at"`
}

type UploadedFile struct {
	FileUUID string `json:"file_uuid" gorm:"Column:file_uuid; Type:varchar; Size:255; primaryKey"`
	// user metadata
	UserID string `json:"username" gorm:"Column:username; Type:varchar; Size:255"`

	// project metatdata
	ProjectID   string `json:"project_id" gorm:"Column:project_id; Type:varchar; Size:255"`
	ProjectPath string `json:"project_path" gorm:"Column:project_path;"`

	Bucket string `json:"bucket" gorm:"Column:bucket; Type:varchar; Size:255"`
	// local path on hard disk or s3 path
	Path     string `json:"path" gorm:"Column:path;"`
	Filename string `json:"name" gorm:"Column:filename; Type:varchar; Size:255"`
	Filesize int64  `json:"size" gorm:"Column:filesize; Type:bigint; Size:255"`

	// file metadata
	UploadSuccess bool `json:"upload_success" gorm:"Column:upload_success; Type:boolean"`
	Deleted       bool `json:"deleted" gorm:"Column:deleted; Type:boolean"`
	Safe          bool `json:"safe" gorm:"Column:safe; Type:boolean"`

	//Tags byte `json:"tags" gorm:"type:jsonb"`

	CreatedAt time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`

	// The time file remains on hard disk waiting for reupload
	UploadExpiredAt time.Time `json:"upload_expired_at" gorm:"Column:upload_expired_at; Type:timestamp"`

	// The time file remains on cloud
	ExpiredAt time.Time `json:"expired_at" gorm:"Column:expired_at; Type:timestamp"`
}

type GeneratedFile struct {
	FileUUID string `json:"file_uuid" gorm:"Column:file_uuid; Type:varchar; Size:255; primaryKey"`

	// run metadata
	UserID       string `json:"username" gorm:"Column:username; Type:varchar; Size:255"`
	RunID        int    `json:"run_id" gorm:"Column:run_id; Type:int; Size:255"`
	RunUUID      string `json:"run_uuid" gorm:"Column:run_uuid; Type:varchar; Size:255"`
	RunName      string `json:"run_name" gorm:"Column:run_name; Type:varchar; Size:255"`
	TaskID       string `json:"task_id" gorm:"Column:task_id; Type:varchar; Size:255"`
	TaskUUID     string `json:"task_uuid" gorm:"Column:task_uuid; Type:varchar; Size:255"`
	TaskName     string `json:"task_name" gorm:"Column:task_name; Type:varchar; Size:255"`
	WorkflowUUID string `json:"workflow_uuid" gorm:"workflow_uuid"`

	ProjectID   string `json:"project_id" gorm:"Column:project_id; Type:varchar; Size:255"`
	ProjectPath string `json:"project_path" gorm:"Column:project_path;"`

	Bucket string `json:"bucket" gorm:"Column:bucket; Type:varchar; Size:255"`
	// local path on hard disk or s3 path
	Path string `json:"path" gorm:"Column:path; Type:varchar; Size:255"`

	Tags byte `json:"tags" gorm:"type:jsonb"`

	Filename string `json:"name" gorm:"Column:filename; Type:varchar; Size:255"`
	Filesize int64  `json:"size" gorm:"Column:filesize; Type:bigint; Size:255"`

	UploadSuccess bool      `json:"upload_success" gorm:"Column:upload_success; Type:boolean"`
	DoneRun       bool      `json:"done_run" gorm:"Column:done_run; Type:boolean"`
	CreatedAt     time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`

	// The time file remains on hard disk waiting for reupload
	UploadExpiredAt time.Time `json:"upload_expired_at" gorm:"Column:upload_expired_at; Type:timestamp"`

	// The time file remains on cloud
	ExpiredAt time.Time `json:"expired_at" gorm:"Column:expired_at; Type:timestamp"`
}

type UserUploadFile struct {
	SampleName     string    `gorm:"sample_name"`
	UserName       string    `gorm:"user_name"`
	LocalPath      string    `gorm:"local_path"`
	UserUploadName string    `gorm:"user_upload_name"`
	Filename       string    `gorm:"filename"`
	Filesize       int64     `gorm:"filesize"`
	Bucket         string    `gorm:"bucket"`
	CreatedAt      time.Time `gorm:"created_at"`
	Total          int64     `gorm:"total"`
}

type GetFileUploadByUserResp struct {
	Samples []*GetUserUploadSampleResp `json:"samples"`
	Total   int64                      `json:"total"`
}

type GetUserUploadSampleResp struct {
	SampleName string                   `json:"sample_name" gorm:"Column:sample_name; Type:varchar; Size:255"`
	UploadTime time.Time                `json:"upload_time" gorm:"Column:upload_time; Type:timestamp"`
	FileList   []*GetUserUploadFileList `json:"file_list"`
}

type GetUserUploadFileList struct {
	FileUUID string `json:"file_uuid" gorm:"Column:path; Type:varchar; Size:255"`
	Path     string `json:"path" gorm:"Column:path; Type:varchar; Size:255"`
	Filename string `json:"name" gorm:"Column:filename; Type:varchar; Size:255"`
	Filesize int64  `json:"size" gorm:"Column:filesize; Type:bigint; Size:255"`
}

type FileHTTPResp struct {
	Path      string    `json:"path" gorm:"Column:path; Type:varchar; Size:255"`
	UserID    string    `json:"username" gorm:"Column:username; Type:varchar; Size:255"`
	RunUUID   string    `json:"run_uuid" gorm:"Column:run_uuid; Type:varchar; Size:255"`
	RunName   string    `json:"run_name" gorm:"Column:run_name; Type:varchar; Size:255"`
	TaskUUID  string    `json:"task_uuid" gorm:"Column:task_uuid; Type:varchar; Size:255"`
	TaskName  string    `json:"task_name" gorm:"Column:task_name; Type:varchar; Size:255"`
	Filename  string    `json:"name" gorm:"Column:filename; Type:varchar; Size:255"`
	Filesize  int64     `json:"size" gorm:"Column:filesize; Type:bigint; Size:255"`
	CreatedAt time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`
	ExpiredAt time.Time `json:"expired_at" gorm:"Column:expired_at; Type:timestamp"`
}

type Bucket struct {
	ID           int       `json:"id" gorm:"Column:id; Type:int"`
	Name         string    `json:"name" gorm:"Column:name; Type:varchar; Size:255"`
	CurrentSize  int64     `json:"current_size" gorm:"Column:current_size; Type:int"`
	CurrentCount int32     `json:"currend_count" gorm:"Column:current_count; Type:int"`
	CreatedAt    time.Time `json:"created_at" gorm:"Column:created_at; Type:timestamp"`
}
