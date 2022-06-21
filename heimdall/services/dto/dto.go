package dto

import (
	"time"

	"github.com/google/uuid"
)

type TransformRequest struct {
	RunIndex int                    `json:"run_index"`
	UserName string                 `json:"username"` // FIXME change all to user name
	Name     string                 `json:"name"`
	Content  string                 `json:"content"`
	Steps    []*WorkflowStep        `json:"steps"`
	Params   map[string]interface{} `json:"workflow_params"`
}

type WorkflowStep struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type TransformResponse struct {
	Error struct {
		Code    uint   `json:"Code"`
		Message string `json:"Message"`
	} `json:"Error"`
	Data Run `json:"Data"`
}

type Run struct {
	WorkflowID  uuid.UUID `json:"workflow_id"`
	RunID       int       `json:"run_id"`
	ProjectID   uuid.UUID `json:"project_id"`   // FUTURE USE
	ProjectPath string    `json:"project_path"` // FUTURE USE
	RunName     string    `json:"run_name"`
	UserName    string    `json:"username"`
	// Status     int     `json:"status"`
	Tasks []*Task `json:"tasks"`
}

type Task struct {
	ID              uuid.UUID         `json:"id"`
	TaskID          string            `json:"task_id"`
	RunID           uuid.UUID         `json:"run_id"`
	ProjectID       uuid.UUID         `json:"project_id"`   // FUTURE USE
	ProjectPath     string            `json:"project_path"` // FUTURE USE
	TaskName        string            `json:"task_name"`
	IsBoundary      bool              `json:"is_boundary"`
	StepName        string            `json:"step_name"`
	UserName        string            `json:"username"`
	Command         string            `json:"command"`
	ParamsWithRegex []*ParamWithRegex `json:"paramwithregex"`
	OutputRegex     []string          `json:"output_regex"`
	Output2ndFiles  []string          `json:"output_2nd_files"`
	ParentTasksID   []string          `json:"parent_tasks_id"`
	ChildrenTasksID []string          `json:"children_tasks_id"`
	OutputLocation  []string          `json:"output_location"`
	DockerImage     []string          `json:"docker_image"`
	QueueLevel      int               `json:"queue_level"`
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

type AuthRequestJSON struct {
	User     AuthRequestJSON_User      `json:"user"`
	Request  *AuthRequestJSON_Request  `json:"request"`
	Requests []AuthRequestJSON_Request `json:"requests"`
}

type AuthRequestJSON_User struct {
	Token string `json:"token"`
	// The Policies field is optional, and if the request provides a token
	// this gets filled in using the Token field.
	Policies  []string `json:"policies,omitempty"`
	Audiences []string `json:"aud,omitempty"`
}

type Constraints = map[string]string

type AuthRequestJSON_Request struct {
	Resource    string      `json:"resource"`
	Action      Action      `json:"action"`
	Constraints Constraints `json:"constraints,omitempty"`
}

type AuthResponse struct {
	Auth bool `json:"auth"`
}

type Action struct {
	Service string `json:"service"`
	Method  string `json:"method"`
}

type Dataset struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Datasets    []Dataset `json:"datasets"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
