package services

import (
	"context"
	"errors"
	"time"

	"workflow/heimdall/core"

	"github.com/google/uuid"

	"workflow/heimdall/repository"
	"workflow/heimdall/repository/entity"
	utilsModel "workflow/workflow-utils/model"
)

//NOTE: Merge UpdateStatusFail, UpdateStatusCanceling, UpdateStatusPausing to 1 func
type iTaskStatusUpdater interface {
	UpdateStatusAuto(ctx context.Context, taskIDs []string, status uint32) (err error)
	UpdateTaskStatusOnly(ctx context.Context, taskIDs []string, state string) (err error)
	UpdateStatusSent(ctx context.Context, runID uuid.UUID) (err error)
	UpdateStatusSuccess(ctx context.Context, taskIDs []string) (err error)
	UpdateStatusFail(ctx context.Context, taskIDs []string) (err error)
	UpdateStatusCanceling(ctx context.Context, taskIDs []string) (err error)
	UpdateStatusCancelled(ctx context.Context, taskIDs []string) (err error)
	UpdateStatusPausing(ctx context.Context, taskIDs []string) (err error) //FIXME: add core.statePausing
	UpdateStatusPaused(ctx context.Context, taskIDs []string) (err error)
}

type taskStatusUpdater struct{}

var iplmTaskUpdater iTaskStatusUpdater = new(taskStatusUpdater)

func GetTaskStatusUpdater() iTaskStatusUpdater {
	return iplmTaskUpdater
}

func (tSU *taskStatusUpdater) UpdateStatusAuto(ctx context.Context, taskIDs []string, status uint32) (err error) {
	switch status {
	case utilsModel.StatusSuccess:
		if err := tSU.UpdateStatusSuccess(ctx, taskIDs); err != nil {
			return err
		}

	case utilsModel.StatusPausing:
		if err := tSU.UpdateStatusPausing(ctx, taskIDs); err != nil {
			return err
		}

	case utilsModel.StatusPaused:
		if err := tSU.UpdateStatusPaused(ctx, taskIDs); err != nil {
			return err
		}

	case utilsModel.StatusCanceling:
		if err := tSU.UpdateStatusCanceling(ctx, taskIDs); err != nil {
			return err
		}

	case utilsModel.StatusCanceled:
		if err := tSU.UpdateStatusCancelled(ctx, taskIDs); err != nil {
			return err
		}

	case utilsModel.StatusFail:
		if err := tSU.UpdateStatusFail(ctx, taskIDs); err != nil {
			return err
		}

	case utilsModel.StatusInqueue:
		if err = tSU.UpdateTaskStatusOnly(ctx, taskIDs, core.StateQueued); err != nil {
			return err
		}

	case utilsModel.StatusInitializing:
		if err = tSU.UpdateTaskStatusOnly(ctx, taskIDs, core.StateInitalizing); err != nil {
			return err
		}

	case utilsModel.StatusInitiated:
		if err = tSU.UpdateTaskStatusOnly(ctx, taskIDs, core.StateRunning); err != nil { //FIXME: missing state
			return err
		}

	case utilsModel.StatusRunning:
		if err = tSU.UpdateTaskStatusOnly(ctx, taskIDs, core.StateRunning); err != nil {
			return err
		}

	default:
		return errors.New("Status code not found")
	}

	return nil
}

func (tSU *taskStatusUpdater) UpdateTaskStatusOnly(ctx context.Context, taskIDs []string, state string) (err error) {
	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	if state == core.StateRunning {
		if _, err := repository.GetTaskDAO().UpdateTaskStartTime(ctx, tx, taskIDs, time.Now()); err != nil {
			return err
		}
	}

	tasks, err := repository.GetTaskDAO().UpdateTaskStatus(ctx, tx, taskIDs, state)
	if err != nil {
		return err
	}

	for taskIndex := range tasks {
		if err = repository.GetTaskDAO().UpdateRunState(ctx, tx, tasks[taskIndex].RunID, core.StateRunning); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (tSU *taskStatusUpdater) UpdateStatusSent(ctx context.Context, runID uuid.UUID) (err error) {
	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	if err := repository.GetTaskDAO().UpdateTaskStatusByRunID(ctx, tx, runID, core.StateInitalizing); err != nil {
		return err
	}

	if err := repository.GetTaskDAO().UpdateRunState(ctx, tx, runID, core.StateRunning); err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (tSU *taskStatusUpdater) UpdateStatusSuccess(ctx context.Context, taskIDs []string) (err error) {
	var (
		taskCount int64 = 0
	)

	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	tasks, err := repository.GetTaskDAO().UpdateTaskStatus(ctx, tx, taskIDs, core.StateComplete)
	if err != nil {
		return err
	}

	for taskIndex := range tasks {
		if err = tx.WithContext(ctx).Model(&entity.TaskEntity{}).
			Where("run_id = ? AND state NOT LIKE ?", tasks[taskIndex].RunID, core.StateComplete).
			Count(&taskCount).Error; err != nil {
			return err
		}
		// the last task (logic task) is incomplete => the number of task not complete == 1
		if taskCount == 1 {
			if err = repository.GetTaskDAO().UpdateRunState(ctx, tx, tasks[taskIndex].RunID, core.StateComplete); err != nil {
				return err
			}
			if err = repository.GetTaskDAO().UpdateTaskStatusByRunID(ctx, tx, tasks[taskIndex].RunID, core.StateComplete); err != nil {
				return err
			}

			logger.Infof("update run state to complete: %s", tasks[taskIndex].RunID)
		}

	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (tSU *taskStatusUpdater) UpdateStatusFail(ctx context.Context, taskIDs []string) (err error) {
	var (
		runIDMap = make(map[uuid.UUID]bool)
	)

	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	tasks, err := repository.GetTaskDAO().UpdateTaskStatus(ctx, tx, taskIDs, core.StateExecutorError)
	if err != nil {
		return err
	}

	for taskIndex := range tasks {
		if _, ok := runIDMap[tasks[taskIndex].RunID]; ok {
			continue
		}

		if err = repository.GetTaskDAO().UpdateTaskStatusByRunID(ctx, tx, tasks[taskIndex].RunID, core.StateExecutorError, core.StateComplete); err != nil {
			return err
		}

		if err = repository.GetTaskDAO().UpdateRunState(ctx, tx, tasks[taskIndex].RunID, core.StateExecutorError); err != nil {
			return err
		}
		logger.Infof("Update Run State to fail: %s", tasks[taskIndex].RunID)

		runIDMap[tasks[taskIndex].RunID] = false
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (tSU *taskStatusUpdater) UpdateStatusCanceling(ctx context.Context, taskIDs []string) (err error) {
	var (
		runIDMap = make(map[uuid.UUID]bool)
	)

	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	tasks, err := repository.GetTaskDAO().UpdateTaskStatus(ctx, tx, taskIDs, core.StateCanceling)
	if err != nil {
		return err
	}

	for taskIndex := range tasks {
		if _, ok := runIDMap[tasks[taskIndex].RunID]; ok {
			continue
		}

		if err = repository.GetTaskDAO().UpdateRunState(ctx, tx, tasks[taskIndex].RunID, core.StateCanceling); err != nil {
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (tSU *taskStatusUpdater) UpdateStatusCancelled(ctx context.Context, taskIDs []string) (err error) {
	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	tasks, err := repository.GetTaskDAO().UpdateTaskStatus(ctx, tx, taskIDs, core.StateCanceled)
	if err != nil {
		return err
	}

	//FIXME: move to DAO
	for taskIndex := range tasks {
		if rows := tx.WithContext(ctx).Model(&entity.TaskEntity{}).Where("run_id = ? AND state LIKE ?", tasks[taskIndex].RunID, core.StateCanceling).RowsAffected; rows == 0 {
			if err = repository.GetTaskDAO().UpdateRunState(ctx, tx, tasks[taskIndex].RunID, core.StateCanceled); err != nil {
				return err
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (tG *taskStatusUpdater) UpdateStatusPausing(ctx context.Context, taskIDs []string) (err error) {
	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	_, err = repository.GetTaskDAO().UpdateTaskStatus(ctx, tx, taskIDs, core.StatePaused)
	if err != nil {
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (tG *taskStatusUpdater) UpdateStatusPaused(ctx context.Context, taskIDs []string) (err error) {
	tx, err := repository.GetTaskDAO().GetTransaction(ctx)
	if err != nil {
		return err
	}

	tasks, err := repository.GetTaskDAO().UpdateTaskStatus(ctx, tx, taskIDs, core.StatePaused)
	if err != nil {
		return err
	}

	for taskIndex := range tasks {
		if rows := tx.WithContext(ctx).Model(&entity.TaskEntity{}).Where("run_id = ? AND state LIKE ?", tasks[taskIndex].RunID, core.StatePaused).RowsAffected; rows == 0 {
			if err = repository.GetTaskDAO().UpdateRunState(ctx, tx, tasks[taskIndex].RunID, core.StatePaused); err != nil {
				return err
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
