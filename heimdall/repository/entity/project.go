package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// ProjectEntity from CWL
type ProjectEntity struct {
	ID               uuid.UUID `gorm:"unique; not null; type:uuid;default:uuid_generate_v4()"`
	Name             string    `gorm:"not null"`
	Description      string    ``
	Summary          string    ``
	Author           string    ``
	AuthResourcePath Ltree     `gorm:"Column:auth_resource_path; Type:ltree; index"`
	CreatedAt        time.Time
	UpdatedAt        sql.NullTime

	Folders   []FolderEntity   `gorm:"foreignKey: ProjectID"`
	Workflows []WorkflowEntity `gorm:"foreignKey:ProjectID"`
}

type FolderEntity struct {
	ID        uuid.UUID `gorm:"unique; not null; type: uuid; default:uuid_generate_v4()"`
	Name      string    `gorm:"not null"`
	Path      string    `gorm:"not null"`
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uuid.UUID
}
