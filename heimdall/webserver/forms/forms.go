// Package forms have structs use to unmarshal the json request
package forms

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ResponseError struct {
	Message    string `json:"msg"`
	StatusCode int    `json:"status_code"`
}

type WorkflowStepDto struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type WorkflowStepForm struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type WorkflowDto struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Class       string                 `json:"class"`
	Steps       []WorkflowStepDto      `json:"steps"`
	Tags        map[string]interface{} `json:"tags,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Author      string                 `json:"author"`
	ProjectId   uuid.UUID              `json:"project_id,omitempty"`
	ProjectName string                 `json:"project_name,omitempty"`
}

type WorkflowsDto struct {
	Workflows     []WorkflowDto `json:"workflows"`
	NextPageToken string        `json:"next_page_token"`
	Total         int64         `json:"total"`
}

type WorkflowForm struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Steps       []WorkflowStepForm     `json:"steps"`
	Class       string                 `json:"class"`
	Tags        map[string]interface{} `json:"tags"`
	ProjectID   uuid.UUID              `json:"project_id"`
}

type ListRunsDto struct {
	Runs          []RunStatusDto `json:"runs"`
	NextPageToken string         `json:"next_page_token"`
	Total         int64          `json:"total"`
}

type RunStatusDto struct {
	ID        uuid.UUID     `json:"id"`
	State     string        `json:"state"`
	User      string        `json:"user"`
	Request   RunRequestDto `json:"request"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
}

type WorkflowRunForm struct {
	WorkflowParams           map[string]interface{} `json:"workflow_params"`
	WorkflowType             string                 `json:"workflow_type"`
	WorkflowTypeVersion      string                 `json:"workflow_type_version"`
	Tags                     map[string]interface{} `json:"tags"`
	WorkflowEngineParameters map[string]interface{} `json:"workflow_engine_parameters"`
	WorkflowURL              string                 `json:"workflow_url"`
	WorkflowAttachments      []string               `json:"workflow_attachments"`
}

type RunDto struct {
	ID       uuid.UUID              `json:"id"`
	Request  RunRequestDto          `json:"request"`
	State    string                 `json:"state"`
	Log      RunLogDto              `json:"run_log"`
	TaskLogs []RunLogDto            `json:"task_logs"`
	Tasks    []SimpleTaskDto        `json:"tasks"`
	Outputs  map[string]interface{} `json:"outputs"`
}

type RunLogDto struct {
	Name      string    `json:"name"`
	Cmd       []string  `json:"cmd"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	StdOut    string    `json:"stdout"`
	StdErr    string    `json:"stderr"`
	ExitCode  string    `json:"exit_code"`
}

type RunRequestDto struct {
	WorkflowParams           map[string]interface{} `json:"workflow_params"`
	WorkflowType             string                 `json:"workflow_type"`
	WorkflowTypeVersion      string                 `json:"workflow_type_version"`
	Tags                     map[string]string      `json:"tags"`
	WorkflowEngineParameters map[string]string      `json:"workflow_engine_parameters"`
	WorkflowURL              string                 `json:"workflow_url"`
}

type ListTasksDto struct {
	Tasks         []TaskDto `json:"tasks"`
	NextPageToken string    `json:"next_page_token"`
	Total         int64     `json:"total"`
}

type TaskDto struct {
	ID           uuid.UUID         `json:"id"`
	State        string            `json:"state"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Inputs       []TaskInputDto    `json:"inputs"`
	Outputs      []TaskOutputDto   `json:"outputs"`
	Resources    TaskResourcesDto  `json:"resources"`
	Executors    []TaskExecutorDto `json:"executors"`
	Volumes      []string          `json:"volumes"`
	Tags         map[string]string `json:"tags"`
	Logs         []TaskLogDto      `json:"logs"`
	CreationTime time.Time         `json:"creation_time"`
	StartedAt    time.Time         `json:"started_at"`
	EndAt        time.Time         `json:"end_at"`
}

type SimpleTaskDto struct {
	ID        uuid.UUID `json:"id"`
	State     string    `json:"state"`
	Name      string    `json:"name"`
	StartedAt time.Time `json:"started_at"`
	EndAt     time.Time `json:"end_at"`
}

type simpleTaskDtoNullTime struct {
	ID        uuid.UUID  `json:"id"`
	State     string     `json:"state"`
	Name      string     `json:"name"`
	StartTime *time.Time `json:"started_at,omitempty"`
	EndTime   *time.Time `json:"end_at,omitempty"`
}

func (t *SimpleTaskDto) MarshalJSON() ([]byte, error) {
	var ret = &simpleTaskDtoNullTime{
		ID:        t.ID,
		State:     t.State,
		Name:      t.Name,
		StartTime: &t.StartedAt,
		EndTime:   &t.EndAt,
	}
	if t.StartedAt.IsZero() {
		ret.StartTime = nil
	}
	if t.EndAt.IsZero() {
		ret.EndTime = nil
	}
	return json.Marshal(ret)
}

type TaskFormDto struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Inputs      []TaskInputDto    `json:"inputs"`
	Outputs     []TaskOutputDto   `json:"outputs"`
	Resources   TaskResourcesDto  `json:"resources"`
	Executors   []TaskExecutorDto `json:"executors"`
	Volumes     []string          `json:"volumes"`
	Tags        map[string]string `json:"tags"`
	RunURL      string            `json:"run_url"`
}

type TaskInputDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Content     string `json:"content"`
}

type TaskOutputDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Path        string `json:"path"`
	Type        string `json:"type"`
}

type TaskResourcesDto struct {
	CPUCores    int    `json:"cpu_cores"`
	Preemptible bool   `json:"preemptible"`
	RAM         int    `json:"ram_gb"`
	Disk        int    `json:"disk_gb"`
	Zones       string `json:"zones"`
}

type TaskLogDto struct {
	Description string                 `json:"description"`
	Logs        []TaskExecutorLogDto   `json:"logs"`
	Metadata    map[string]string      `json:"metadata"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	Outputs     []TaskOutputFileLogDto `json:"outputs"`
	SystemLogs  []string               `json:"system_logs"`
}

type TaskExecutorLogDto struct {
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	StdOut      string    `json:"stdout"`
	StdErr      string    `json:"stderr"`
	ExitCode    int       `json:"exit_code"`
}

type TaskExecutorDto struct {
	Image   string            `json:"image"`
	Command []string          `json:"command"`
	WorkDir string            `json:"workdir"`
	StdIn   string            `json:"stdin"`
	StdOut  string            `json:"stdout"`
	StdErr  string            `json:"stderr"`
	Env     map[string]string `json:"env"`
}

type TaskOutputFileLogDto struct {
	URL       string `json:"url"`
	Path      string `json:"path"`
	SizeBytes string `json:"size_bytes"`
}

type ServiceInfoDto struct {
	ID               uuid.UUID    `json:"id"`
	Name             string       `json:"name"`
	Type             ServiceType  `json:"type"`
	Description      string       `json:"description"`
	Organization     Organization `json:"organization"`
	ContactURL       string       `json:"contactUrl"`
	DocumentationURL string       `json:"documentationUrl"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	Environment      string       `json:"environment"`
	Version          string       `json:"version"`
	Storage          []string     `json:"storage"`
}

type ServiceType struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
}

type Organization struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ResponseRunTransformer struct {
	Error errorResponseRunTransformer
	Data  interface{}
}

type errorResponseRunTransformer struct {
	Code    int
	Message string // need to be empty if success
}

type ProjectDto struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Author      string      `json:"author"`
	Folders     []FolderDto `json:"folders"`
}

type ProjectForm struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Summary     string   `json:"summary"`
	Folders     []string `json:"folders"`
}

type ProjectsDto struct {
	ShareProjects []ProjectDto `json:"share_projects"`
	Projects      []ProjectDto `json:"projects"`
	NextPageToken string       `json:"next_page_token"`
	Total         int64        `json:"total"`
}

type FolderCreate struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type FolderUpdate struct {
	ID uuid.UUID `json:"id"`
	FolderCreate
}

type FolderDto struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
