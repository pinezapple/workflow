package model

type Arguments struct{
	ToolId uint32 `gorm:"column:tool_of_argument" json:"tool_id"`
	ArgumentId string `gorm:"column:argument_id; primaryKey" json:"argument_id"`
	Argument string `gorm:"column:argument" json:"argument"`
	DefaultValue string `gorm:"column:default_value" json:"default_value"`
	Description string `gorm:"column:description" json:"description"`
	Tool Tool `gorm:"foreignKey:tool_of_argument;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}