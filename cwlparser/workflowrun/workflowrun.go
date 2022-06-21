package workflowrun

type Run struct {
	RunIndex int     `json:"run_index"`
	RunName  string  `json:"run_name"`
	UserName string  `json:"username"`
	Tasks    []*Task `json:"tasks"`
}

type Task struct {
	TaskID          string            `json:"task_id"`
	TaskName        string            `json:"task_name"`
	StepID          string            `json:"-"`
	UserName        string            `json:"username"`
	Command         string            `json:"command"`
	ScatterMethod   string            `json:"scatter_method"`
	ParamsWithRegex []*ParamWithRegex `json:"paramwithregex"`
	OutputRegex     []string          `json:"output_regex"`
	Output2ndFiles  []string          `json:"output_2nd_files"`
	ParentTasksID   []string          `json:"parent_tasks_id"`
	ChildrenTasksID []string          `json:"children_tasks_id"`
	OutputLocation  []string          `json:"output_location"`
	DockerImage     []string          `json:"docker_image"`
	IsBoundary      bool              `json:"is_boundary"`
	RunIndex        int               `json:"run_index"`
	QueueLevel      int               `json:"queue_level"`
	Status          int               `json:"status"`
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

type NewRun struct {
	WorkflowID   string
	WorkflowUUID string
	Username     string
	State        string
	Tasks        []struct {
		TaskID     string
		TaskUUID   string
		TaskName   string
		IsBoundary bool

		StepID   string
		StepName string

		DockerImage string

		State string

		ParentIDs     []string
		ParentUUIDs   []string
		ChildrenIDs   []string
		ChildrenUUIDs []string

		ScatterMethod     string
		ScatterParamNames []string

		Command []string // joins + " "

		Params []struct {
			Name           string
			From           string
			Prefix         string
			IsScatter      bool
			SecondaryFiles []string
			Patterns       []string
			Values         []struct {
				IsFile    bool
				FileSizes []int64
				Values    []string
			}
		}

		Outputs []struct {
			Name           string
			Patterns       []string
			SecondaryFiles []string
		}
	}
}

type NewTask struct {
	TaskID     string
	TaskUUID   string
	TaskName   string
	IsBoundary bool

	StepID   string
	StepName string

	DockerImage string

	State string

	ParentIDs     []string
	ParentUUIDs   []string
	ChildrenIDs   []string
	ChildrenUUIDs []string

	ScatterMethod string
	ScatterParam  []string

	Command []string // joins + " "

	Params []Param

	Outputs []Output
}

type Param struct {
	Name           string
	From           string
	Prefix         string
	IsScatter      bool
	SecondaryFiles []string
	Patterns       []string
	Values         []Value
}

type Value struct {
	IsFile    bool
	FileSizes []int64
	Values    []string
}

type Output struct {
	Name           string
	Patterns       []string
	SecondaryFiles []string
}
