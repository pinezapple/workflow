package gormdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/vfluxus/heimdall/core"
	"github.com/vfluxus/heimdall/repository/entity"
	"github.com/vfluxus/heimdall/utils"
	"github.com/vfluxus/workflow-utils/model"
)

type TaskGorm struct {
	sync.Mutex
}

// Get Transaction
func (tG *TaskGorm) GetTransaction(ctx context.Context) (tx *gorm.DB, err error) {
	tG.Lock()
	tx = gDB.Begin()
	tG.Unlock()

	return tx, nil
}

func (tG *TaskGorm) UpdateTaskStatus(ctx context.Context, tx *gorm.DB, taskIDs []string, state string) (tasks []*entity.TaskEntity, err error) {
	if tx == nil {
		tx = gDB
	}

	for taskIDIndex := range taskIDs {
		task := new(entity.TaskEntity)
		var values entity.TaskEntity
		if state == core.StateCanceled || state == core.StateComplete || state == core.StateExecutorError || state == core.StateSystemError {
			values = entity.TaskEntity{State: state, EndTime: time.Now()}
		} else {
			values = entity.TaskEntity{State: state}
		}
		if err := tx.WithContext(ctx).Model(&entity.TaskEntity{}).Where("task_id = ? ", taskIDs[taskIDIndex]).Take(task).Updates(values).Error; err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tG *TaskGorm) UpdateTaskStartTime(ctx context.Context, tx *gorm.DB, taskIDs []string, starttime time.Time) (tasks []*entity.TaskEntity, err error) {
	if tx == nil {
		tx = gDB
	}

	for taskIDIndex := range taskIDs {
		task := new(entity.TaskEntity)
		values := entity.TaskEntity{StartedTime: starttime}
		if err := tx.WithContext(ctx).Model(&entity.TaskEntity{}).Where("task_id = ? ", taskIDs[taskIDIndex]).Take(task).Updates(values).Error; err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tG *TaskGorm) UpdateRunState(ctx context.Context, tx *gorm.DB, runID uuid.UUID, state string) (err error) {
	var values entity.RunEntity
	if state == core.StateCanceled || state == core.StateComplete || state == core.StateExecutorError || state == core.StateSystemError {
		values = entity.RunEntity{State: state, EndTime: sql.NullTime{Time: time.Now(), Valid: true}}
	} else {
		values = entity.RunEntity{State: state}
	}
	if err := tx.WithContext(ctx).Where("id = ?", runID).Take(&entity.RunEntity{}).Updates(values).Error; err != nil {
		return err
	}

	return nil
}

func (tG *TaskGorm) UpdateTaskStatusByRunID(ctx context.Context, tx *gorm.DB, runID uuid.UUID, state string, excludeState ...string) (err error) {
	var (
		stmt strings.Builder
	)

	if tx == nil {
		tx = gDB
	}

	// build stmt
	for i := range excludeState {
		_, err = stmt.WriteString("AND state NOT LIKE '" + excludeState[i] + "' ")
		if err != nil {
			return err
		}
	}

	if err = tx.WithContext(ctx).Model(&entity.TaskEntity{}).Where("run_id = ? "+stmt.String(), runID).Take(&entity.TaskEntity{}).Update("state", state).Error; err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- UPDATE INFO ----------------------------------------------------------

func (tG *TaskGorm) UpdateOutput(ctx context.Context, taskID string, output []byte) (err error) {
	var (
		updateTask = &entity.TaskEntity{
			TaskID:  taskID,
			Outputs: output,
		}
	)

	if err = gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Where("task_id = ?", taskID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
		return err
	}

	return nil
}

func (tG *TaskGorm) UpdateOutputLocation(ctx context.Context, taskID string, output []string) (err error) {
	var (
		checkTask  = new(entity.TaskEntity)
		updateTask = &entity.TaskEntity{
			TaskID:         taskID,
			OutputLocation: output,
		}
	)
	if err := gDB.WithContext(ctx).Where("task_id = ?", taskID).Take(checkTask).Updates(updateTask).Error; err != nil {
		return err
	}

	return nil
}

func (tG *TaskGorm) UpdateParamWithRegex(ctx context.Context, taskID string, paramWithRegex []byte) (err error) {
	var (
		updateTask = &entity.TaskEntity{
			TaskID:          taskID,
			ParamsWithRegex: paramWithRegex,
		}
	)

	if err = gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Where("task_id = ?", taskID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
		return err
	}

	return nil
}

//TODO: Update File size field to Task
func (tG *TaskGorm) UpdateFileSize(ctx context.Context, taskID string, fileSize []int) (err error) {
	var (
		updateTask = &entity.TaskEntity{
			TaskID: taskID,
			// FileSize: fileSize,
		}
	)

	if err = gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Where("task_id = ?", taskID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
		return err
	}

	return nil
}

func (tG *TaskGorm) UpdateDoneTask(ctx context.Context, taskID string, outputFileName []string, outputFileSize []int64, childrenInputFile, childrenInputFileName []string, childrenInputFilesize []int64) (err error) {
	return gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task = entity.TaskEntity{}
		err = tx.Model(&entity.TaskEntity{}).Where(&entity.TaskEntity{TaskID: taskID}).Find(&task).Error
		if task.State == core.StateExecutorError {
			return nil
		}
		// Update task
		var (
			updateTask = &entity.TaskEntity{
				TaskID:         taskID,
				OutputLocation: outputFileName,
				OutputFilesize: outputFileSize,
				State:          core.StateComplete,
			}
		)
		if err := tx.Where("id = ?", task.ID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
			return err
		}

		var childrenTasks []entity.TaskEntity
		err = tx.Raw("SELECT FROM task_entities WHERE id IN (SELECT children_task_id FROM task_entities WHERE id = ?) AND parent_done_count <> 0", task.ID).Scan(&childrenTasks).Error
		for i := 0; i < len(childrenTasks); i++ {
			// filling child param with regexes
			var childParam []*model.ParamWithRegex
			err := json.Unmarshal(childrenTasks[i].ParamsWithRegex, &childParam)
			if err != nil {
				return err
			}
			utils.FillInput(childParam, childrenInputFileName, childrenInputFile, childrenInputFilesize)
			childParamByte, err := json.Marshal(childParam)
			if err != nil {
				return err
			}

			var (
				updateTask = &entity.TaskEntity{
					ID:               childrenTasks[i].ID,
					ParamsWithRegex:  childParamByte,
					ParentsDoneCount: childrenTasks[i].ParentsDoneCount - 1,
				}
			)

			// update to db
			if err = tx.Model(&entity.TaskEntity{}).Where("id= ?", childrenTasks[i].ID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
				return err
			}
		}

		// return nil will commit the whole transaction
		return nil
	})
}

func (tG *TaskGorm) UpdateFailTask(ctx context.Context, taskID string) (err error) {
	return gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task = entity.TaskEntity{}
		err = tx.Model(&entity.TaskEntity{}).Where(&entity.TaskEntity{TaskID: taskID}).Find(&task).Error
		if task.State == core.StateExecutorError {
			return nil
		}
		if err = tx.WithContext(ctx).Model(&entity.TaskEntity{}).Where("run_id = ? ", task.RunID).Take(&entity.TaskEntity{}).Update("state", core.StateExecutorError).Error; err != nil {
			return err
		}

		var values entity.RunEntity
		values = entity.RunEntity{State: core.StateExecutorError, EndTime: sql.NullTime{Time: time.Now(), Valid: true}}
		if err := tx.WithContext(ctx).Where("id = ?", task.RunID).Take(&entity.RunEntity{}).Updates(values).Error; err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})
}

func (tG *TaskGorm) GetTaskByTaskID(ctx context.Context, id string) (task entity.TaskEntity, err error) {
	err = gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Where(&entity.TaskEntity{TaskID: id}).Find(&task).Error
	return
}

func (tG *TaskGorm) GetReadyChildrenTaskByTaskID(ctx context.Context, taskID string) (task []entity.TaskEntity, err error) {
	err = gDB.WithContext(ctx).Raw("SELECT FROM task_entities WHERE id IN (SELECT children_task_id FROM task_entities WHERE task_id = ?) AND parent_done_count = 0", taskID).Scan(&task).Error
	return
}

func (tG *TaskGorm) UpdateTaskState(ctx context.Context, taskID string, state string) (err error) {
	return gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task = entity.TaskEntity{}
		err = tx.Model(&entity.TaskEntity{}).Where(&entity.TaskEntity{TaskID: taskID}).Find(&task).Error
		if task.State == core.StateExecutorError {
			return utils.ErrTaskFailed
		}
		// Update task
		var (
			updateTask = &entity.TaskEntity{
				TaskID: taskID,
				State:  state,
			}
		)
		if err := tx.Where("task_id = ?", task.TaskID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
			return err
		}
		return nil
	})
}

func (tG *TaskGorm) UpdateStartTime(ctx context.Context, taskID string, starttime time.Time) (err error) {
	task := new(entity.TaskEntity)
	values := entity.TaskEntity{StartedTime: starttime}
	if err := gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Where("task_id = ? ", taskID).Take(task).Updates(values).Error; err != nil {
		return err
	}

	return nil
}
