package gormdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"workflow/heimdall/core"
	"workflow/heimdall/repository/entity"
	"workflow/heimdall/utils"
	"workflow/workflow-utils/model"
)

type TaskGorm struct {
}

func (tG *TaskGorm) UpdateDoneTask(ctx context.Context, taskID string, outputFileName []string, outputFileSize []int64, childrenInputFile, childrenInputFileName []string, childrenInputFilesize []int64) (err error) {
	return gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task = entity.TaskEntity{}
		err = tx.Model(&entity.TaskEntity{}).Where(&entity.TaskEntity{TaskID: taskID}).Find(&task).Error
		if task.State == core.StateExecutorError || task.State == core.StateComplete {
			return utils.ErrTaskDone
		}

		// Update task
		var (
			updateTask = &entity.TaskEntity{
				OutputLocation: outputFileName,
				OutputFilesize: outputFileSize,
				State:          core.StateComplete,
				EndTime:        time.Now(),
			}
		)
		if err := tx.Where("id = ?", task.ID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
			return err
		}

		var childrenTasks []entity.TaskEntity
		err = tx.Raw("SELECT * FROM task_entities WHERE task_id = ANY((SELECT children_tasks_id FROM task_entities WHERE task_id = ?)::text[]) AND parent_done_count <> 1", task.TaskID).Scan(&childrenTasks).Error

		for i := 0; i < len(childrenTasks); i++ {
			// This is the final task of the run
			if childrenTasks[i].IsBoundary && childrenTasks[i].ParentDoneCount == 2 {
				// Update Done Run
				var values entity.RunEntity
				values = entity.RunEntity{State: core.StateComplete, EndTime: sql.NullTime{Time: time.Now(), Valid: true}}
				if err := tx.WithContext(ctx).Where("id = ?", task.RunID.String()).Take(&entity.RunEntity{}).Updates(values).Error; err != nil {
					return err
				}
			}

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
			parentDoneCount := childrenTasks[i].ParentDoneCount - 1

			updateTask = &entity.TaskEntity{
				ParamsWithRegex: childParamByte,
				ParentDoneCount: parentDoneCount,
			}

			// update to db
			if err = tx.Where("id = ?", childrenTasks[i].ID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
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
		if task.State == core.StateExecutorError || task.State == core.StateComplete {
			return nil
		}
		// Update task fail
		var (
			updateTask = &entity.TaskEntity{
				State:   core.StateExecutorError,
				EndTime: time.Now(),
			}
		)
		if err := tx.Where("id = ?", task.ID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
			return err
		}

		if err := tx.Exec("UPDATE task_entities SET state = ?, started_time = ?, end_time = ?  WHERE run_id = ? AND state = ?", core.StateExecutorError, time.Now(), time.Now(), task.RunID.String(), core.StateUnknown).Error; err != nil {
			return err
		}

		/*
			if err = tx.WithContext(ctx).Model(&entity.TaskEntity{}).Where("run_id = ? ", task.RunID.String()).Take(&entity.TaskEntity{}).Update("state", core.StateExecutorError).Error; err != nil {
				return err
			}
		*/

		var values entity.RunEntity
		values = entity.RunEntity{State: core.StateExecutorError, EndTime: sql.NullTime{Time: time.Now(), Valid: true}}
		if err := tx.WithContext(ctx).Where("id = ?", task.RunID.String()).Take(&entity.RunEntity{}).Updates(values).Error; err != nil {
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
	err = gDB.WithContext(ctx).Raw("SELECT * FROM task_entities WHERE task_id = ANY((SELECT children_tasks_id FROM task_entities WHERE task_id = ?)::text[]) AND parent_done_count = 1 AND state = ?", taskID, core.StateUnknown).Scan(&task).Error
	return
}

func (tG *TaskGorm) GetChildrenTaskByTaskID(ctx context.Context, taskID string) (task []entity.TaskEntity, err error) {
	err = gDB.WithContext(ctx).Raw("SELECT * FROM task_entities WHERE task_id = ANY((SELECT children_tasks_id FROM task_entities WHERE task_id = ?)::text[]) AND state = ?", taskID, core.StateUnknown).Scan(&task).Error
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

		if state == core.StateQueued {
			if err := tx.Exec("UPDATE task_entities SET parent_done_count = ? WHERE id = ?", 1, task.ID).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("task_id = ?", task.TaskID).Take(&entity.TaskEntity{}).Updates(updateTask).Error; err != nil {
			return err
		}
		return nil
	})
}

func (tG *TaskGorm) UpdateStartTime(ctx context.Context, taskID string, starttime time.Time) (err error) {
	task := new(entity.TaskEntity)
	values := entity.TaskEntity{StartedTime: starttime, State: core.StateRunning}
	if err := gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Where("task_id = ? ", taskID).Take(task).Updates(values).Error; err != nil {
		return err
	}

	return nil
}
