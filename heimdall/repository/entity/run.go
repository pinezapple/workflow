package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// RunEntity contains run information received from client and executor, scheduler
type RunEntity struct {
	ID          uuid.UUID `gorm:"unique; not null;type:uuid; default: uuid_generate_v4()"`
	ProjectID   uuid.UUID `gorm:"type:uuid"`
	ProjectPath string
	Description string
	Tags        []byte `gorm:"type:jsonb"`
	Request     []byte `gorm:"type:jsonb"`
	WorkflowID  uuid.UUID
	UserName    string
	State       string
	RunLog      []byte `gorm:"type:jsonb"`
	RunIndex    int    `gorm:"not null; autoIncrement"`
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
	EndTime     sql.NullTime

	Tasks []TaskEntity `gorm:"foreignKey:RunID; constraint:OnDelete:CASCADE"`
}

// TaskEntity contains task information received from executor
type TaskEntity struct {
	ID     uuid.UUID `gorm:"not null; unique;type:uuid;default: uuid_generate_v4()"`
	TaskID string    `gorm:"not null"`
	Name   string    `json:"task_name"`

	ProjectID   uuid.UUID `gorm:"type:uuid"`
	ProjectPath string    `json:"project_path" gorm:"Column:project_path;"`

	StepName string

	RunID    uuid.UUID `json:"run_id" gorm:"column:run_id"`
	RunIndex int

	UserName    string `json:"username" gorm:"column:username"`
	Description string
	IsBoundary  bool `json:"is_boundary" gorm:"column:is_boundary"`

	OutputRegex     pq.StringArray `gorm:"type:text[]"`
	Output2ndFiles  pq.StringArray `gorm:"type:text[]"`
	OutputLocation  pq.StringArray `gorm:"type:text[]"`
	OutputFilesize  pq.Int64Array  `json:"output_file_size" gorm:"type:bigint[];column:output_file_size"`
	ParamsWithRegex []byte         `gorm:"type:jsonb"`

	ParentTasksID   pq.StringArray `gorm:"type:text[]"`
	ChildrenTasksID pq.StringArray `gorm:"type:text[]"`

	Command     pq.StringArray `gorm:"type:text[]"`
	RealCommand string
	DockerImage pq.StringArray `gorm:"type:text[]"`
	Inputs      []byte         `gorm:"type:jsonb"`
	Outputs     []byte         `gorm:"type:jsonb"`
	Resource    []byte         `gorm:"type:jsonb"`
	Executors   []byte         `gorm:"type:jsonb"`
	Logs        []byte         `gorm:"type:jsonb"`

	State           string
	ParentDoneCount int `json:"parent_done_count" gorm:"column:parent_done_count"`

	StartedTime time.Time // start time in executor
	EndTime     time.Time // end time in executor

	CreatedAt time.Time
	UpdatedAt time.Time
}
