package model

type Role struct{
	RoleId uint8 `gorm:"column:role_id; primaryKey; autoIncrement:false" json:"role_id"`
	FileLimit uint8 `gorm:"column:file_limit" json:"file_limit"`
	CpuLimit uint8 `gorm:"column:cpu_limit" json:"cpu_limit"`
	RamLimit uint8 `gorm:"column:ram_limit" json:"ram_limit"`
}
