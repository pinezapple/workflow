package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"workflow/heimdall/repository/entity"

	"workflow/heimdall/repository/gormdb"
	"gorm.io/gorm"
)

type iTaskDAO interface {
	// Get DB Transaction
	GetTransaction(ctx context.Context) (tx *gorm.DB, err error)

	//Update status
	UpdateTaskStatus(ctx context.Context, tx *gorm.DB, taskIDs []string, state string) (tasks []*entity.TaskEntity, err error)
	UpdateTaskStartTime(ctx context.Context, tx *gorm.DB, taskIDs []string, starttime time.Time) (tasks []*entity.TaskEntity, err error)
	UpdateTaskStatusByRunID(ctx context.Context, tx *gorm.DB, runID uuid.UUID, state string, excludeState ...string) (err error)

	// Update Run
	UpdateRunState(ctx context.Context, tx *gorm.DB, runID uuid.UUID, state string) (err error)

	// Update task info
	UpdateOutput(ctx context.Context, taskID string, output []byte) (err error)
	UpdateOutputLocation(ctx context.Context, taskID string, output []string) (err error)
	UpdateParamWithRegex(ctx context.Context, taskID string, paramWithRegex []byte) (err error)
	UpdateFileSize(ctx context.Context, taskID string, fileSize []int) (err error)

	GetTaskByTaskID(ctx context.Context, id string) (task entity.TaskEntity, err error)
	GetReadyChildrenTaskByTaskID(ctx context.Context, taskID string) (task []entity.TaskEntity, err error)
	UpdateDoneTask(ctx context.Context, taskID string, outputFileName []string, outputFileSize []int64, childrenInputFile, childrenInputFileName []string, childrenInputFilesize []int64) (err error)
	UpdateFailTask(ctx context.Context, taskID string) (err error)
	UpdateTaskState(ctx context.Context, taskID string, state string) (err error)
	UpdateStartTime(ctx context.Context, taskID string, starttime time.Time) (err error)
}

var (
	iplmTaskDAO iTaskDAO = &gormdb.TaskGorm{}
)

func GetTaskDAO() iTaskDAO {
	return iplmTaskDAO
}
