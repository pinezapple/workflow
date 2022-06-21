package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// WorkflowEntity from CWL
type WorkflowEntity struct {
	ID          uuid.UUID `gorm:"unique; not null; type:uuid; default: uuid_generate_v4()"`
	Name        string    `gorm:"not null"`
	Description string
	Content     string
	Steps       []WorkflowStepEntity `gorm:"foreignKey:WorkflowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tags        []byte               `gorm:"type:jsonb"`
	Class       string               `gorm:"default:workflow"` // maybe workflow or tool
	Author      string
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
	// DeletedAt   gorm.DeletedAt `gorm:"index"`

	ProjectID uuid.UUID
	Project   ProjectEntity `grom:"foreignKey:ProjectID"`
}

// WorkflowStepEntity from CWL
type WorkflowStepEntity struct {
	ID         uuid.UUID `gorm:"unique; not null; type:uuid; default:uuid_generate_v4()"`
	Name       string    `gorm:"not null"`
	Content    string
	WorkflowID uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
	// DeletedAt  gorm.DeletedAt `gorm:"index"`
	Workflow WorkflowEntity `gorm:"foreignKey:WorkflowID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
