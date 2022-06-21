package dto

type SchedRun struct {
	WorkflowID  string `json:"workflow_id"`
	RunID       int    `json:"run_id"`
	RunName     string `json:"run_name"`
	ProjectID   string `json:"project_id"`
	ProjectPath string `json:"project_path"`
	UserName    string `json:"username"`
	// Status     int     `json:"status"`
	Tasks []*SchedTask `json:"tasks"`
}

type SchedTask struct {
	ID              string            `json:"id"`
	TaskID          string            `json:"task_id"`
	ProjectID       string            `json:"project_id"`
	ProjectPath     string            `json:"project_path"`
	RunID           int               `json:"run_id"`
	RunUUID         string            `json:"run_uuid"`
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
