package model

type Task struct {
	TaskId       string `gorm:"column:task_id; primaryKey" json:"task_id"`
	DagId        string `gorm:"column:dag_of_task" json:"dag_id"`
	DAG          DAG    `gorm:"foreignKey:dag_of_task;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Taskname     string `gorm:"column:taskname" json:"taskname"`
	ToolId       uint32 `gorm:"column:tool_of_task" json:"tool_id"`
	Tool         Tool   `gorm:"foreignKey:tool_of_task;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Requirement  string `gorm:"column:requirement" json:"requirement"`
	Command      string `gorm:"column:command" json:"command"`
	Param        string `gorm:"column:param" json:"param"`
	Value        string `gorm:"column:value" json:"value"`
	ParentTaskId string `gorm:"column:parent_task_id" json:"parent_task_id"`
}

type TaskDTO struct {
	TaskID          string            `gorm:"column:task_id; primary_key" json:"task_id"`
	TaskUUID        string            `gorm:"column:task_uuid;" json:"task_uuid"`
	Command         string            `json:"command" gorm:"column:command"`
	ParamsWithRegex []*ParamWithRegex `json:"paramwithregex"`
	OutputRegex     []string          `json:"output_regex"`
	Output2ndFiles  []string          `json:"output_2nd_files"`
	DockerImage     []string          `json:"docker_image" gorm:"type:text[];column:docker_image"`
}
