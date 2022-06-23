package repository

import (
	"context"
	"time"

	"workflow/heimdall/repository/entity"

	"workflow/heimdall/repository/gormdb"
)

type iTaskDAO interface {
	GetTaskByTaskID(ctx context.Context, id string) (task entity.TaskEntity, err error)
	GetReadyChildrenTaskByTaskID(ctx context.Context, taskID string) (task []entity.TaskEntity, err error)
	GetChildrenTaskByTaskID(ctx context.Context, taskID string) (task []entity.TaskEntity, err error)
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
