// Package repository provides the methods to interact with database
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vfluxus/heimdall/repository/entity"
	"github.com/vfluxus/heimdall/repository/gormdb"
)

// ------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- DECLARE THE INTERFACES ----------------------------------------------------------
type dao interface {
	AutoMigrate(ctx context.Context) error

	InitDBConnection(dsn string) error

	GetWorkflow(ctx context.Context, id uuid.UUID) (entity.WorkflowEntity, error)
	GetWorkflows(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) ([]entity.WorkflowEntity, int64, error)
	CreateWorkflow(ctx context.Context, workflow *entity.WorkflowEntity) error
	UpdateWorkflow(ctx context.Context, id uuid.UUID, workflow *entity.WorkflowEntity) error
	DeleteWorkflow(ctx context.Context, id uuid.UUID) error
	GetRunsByWorkflowID(ctx context.Context, username string, id uuid.UUID, pageSize int, pageToken int, filter map[string][]string) ([]entity.RunEntity, int64, error)

	GetRun(ctx context.Context, runID uuid.UUID) (entity.RunEntity, error)
	GetRuns(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) ([]entity.RunEntity, int64, error)
	CreateRun(ctx context.Context, run *entity.RunEntity) error
	DeleteRun(ctx context.Context, run *entity.RunEntity) error
	UpdateRun(ctx context.Context, run *entity.RunEntity) error

	GetTask(ctx context.Context, id uuid.UUID) (entity.TaskEntity, error)
	GetTasks(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) ([]entity.TaskEntity, int64, error)
	CreateTask(ctx context.Context, task *entity.TaskEntity) error

	GetProject(ctx context.Context, id uuid.UUID) (entity.ProjectEntity, error)
	GetProjects(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) ([]entity.ProjectEntity, int64, error)
	GetProjectsFromAuth(ctx context.Context, pageSize int, pageToken int, authPath []string) ([]entity.ProjectEntity, int64, error)
	GetAllProjectsFromAuth(ctx context.Context, authPaths []string) (projects []entity.ProjectEntity, err error)
	CreateProject(ctx context.Context, project *entity.ProjectEntity) error
	UpdateProject(ctx context.Context, id uuid.UUID, project *entity.ProjectEntity) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
	GetWorkflowsByProjectID(ctx context.Context, id uuid.UUID, pageSize int, pageToken int, filter map[string][]string) ([]entity.WorkflowEntity, int64, error)
	AddProjectFolder(ctx context.Context, folder *entity.FolderEntity) error
	UpdateProjectFolder(ctx context.Context, folder *entity.FolderEntity) error
	DeleteProjectFolder(ctx context.Context, folder_id uuid.UUID) error
}

var d dao = gormdb.GetGormDB()

func GetDAO() dao {
	return d
}
