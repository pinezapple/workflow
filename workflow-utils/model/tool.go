package model

type Tool struct{
	ToolId uint32 `gorm:"column:tool_id; primaryKey; autoIncrement:false" json:"tool_id"`
	Command string `gorm:"column:command" json:"command"`
	Version string `gorm:"column:version" json:"version"`
	Description string `gorm:"column:description" json:"description"`
}
