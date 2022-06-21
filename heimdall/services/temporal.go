package services

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"github.com/vfluxus/heimdall/core"
	"github.com/vfluxus/heimdall/repository"
	"github.com/vfluxus/heimdall/repository/entity"
	"github.com/vfluxus/heimdall/utils"
	"github.com/vfluxus/workflow-utils/model"
)

const ()

type HeimdallTemporal struct {
	tempCli client.Client
	worker  worker.Worker
}

type ExecuteTaskParam struct {
	Task entity.TaskEntity
}

type UpdateTaskStatusParam struct {
	TaskID string
	Status string
}

type UpdateTaskStartTimeParam struct {
	TaskID    string
	TimeStamp time.Time
}

// Service implementation
func CreateHeimdallTemporal(cli client.Client) *HeimdallTemporal {
	return &HeimdallTemporal{
		tempCli: cli,
	}
}

func (e *HeimdallTemporal) RegisterWorker() (err error) {
	// TODO: fix this after you have config
	workerOptions := worker.Options{
		MaxConcurrentWorkflowTaskExecutionSize: 1000,
	}
	// TODO: add task queue name
	e.worker = worker.New(e.tempCli, "your_task_queue_name", workerOptions)

	// register workflow
	e.worker.RegisterWorkflowWithOptions(e.ExecuteTaskWf, workflow.RegisterOptions{Name: model.ExecuteTaskWfName})

	// register activity
	e.worker.RegisterActivityWithOptions(e.UpdateTaskStatusAct, activity.RegisterOptions{Name: model.UpdateTaskStatusActName})
	e.worker.RegisterActivityWithOptions(e.UpdateTaskSuccessAct, activity.RegisterOptions{Name: model.UpdateTaskSuccessActName})
	e.worker.RegisterActivityWithOptions(e.UpdateTaskFailAct, activity.RegisterOptions{Name: model.UpdateTaskFailActName})

	// TODO: add LOGGGG
	if err := e.worker.Start(); err != nil {
		return err
	}
	return nil
}

func (u *HeimdallTemporal) ExecuteTaskWf(ctx workflow.Context, param ExecuteTaskParam) (err error) {
	log := workflow.GetLogger(ctx)

	// STEP 1: update task status to inqueue
	updateStatusParam := UpdateTaskStatusParam{
		TaskID: param.Task.TaskID,
		Status: core.StateQueued,
	}
	future := workflow.ExecuteActivity(ctx, model.UpdateTaskStatusActName, updateStatusParam)
	if err = future.Get(ctx, nil); err != nil {
		if err == utils.ErrTaskFailed {
			return nil
		}
		log.Error(err.Error())
		return
	}

	// STEP 2: execute task on Executor
	taskDTO, err := transformToTaskDTO(param.Task)
	if err != nil {
		return err
	}
	var req = model.ExecuteTaskParam{
		Task: taskDTO,
	}
	var resp = model.ExecuteTaskResult{}
	future = workflow.ExecuteActivity(ctx, model.ExecuteTaskActName, req)
	if err = future.Get(ctx, &resp); err != nil {
		// if fail to execute job on k8s, fail this run
		future = workflow.ExecuteActivity(ctx, model.UpdateTaskFailActName, model.UpdateTaskFailParam{TaskID: param.Task.TaskID})
		if err = future.Get(ctx, nil); err != nil {
			if err == utils.ErrTaskFailed {
				return nil
			}
			log.Error(err.Error())
			return
		}

		log.Error(err.Error())
		return
	}

	// STEP 3: update task start time
	future = workflow.ExecuteActivity(ctx, model.UpdateTaskStartTimeActName, UpdateTaskStartTimeParam{TaskID: param.Task.TaskID, TimeStamp: resp.TimeStamp})
	if err = future.Get(ctx, nil); err != nil {
		log.Error(err.Error())
		return
	}

	return nil
}

func (u *HeimdallTemporal) UpdateTaskStartTimeAct(ctx context.Context, param UpdateTaskStartTimeParam) error {
	taskDAO := repository.GetTaskDAO()
	return taskDAO.UpdateStartTime(ctx, param.TaskID, param.TimeStamp)
}

func (u *HeimdallTemporal) UpdateTaskStatusAct(ctx context.Context, param UpdateTaskStatusParam) error {
	taskDAO := repository.GetTaskDAO()
	return taskDAO.UpdateTaskState(ctx, param.TaskID, param.Status)
}

func (u *HeimdallTemporal) UpdateTaskSuccessAct(ctx context.Context, param model.UpdateTaskSuccessParam) (res model.UpdateTaskSuccessResult, err error) {
	taskDAO := repository.GetTaskDAO()
	var outputFileName, filePathToSave, fileNameToSave []string
	var outputFileSize, fileSizeToSave []int64

	// STEP 1: Update output location + status + children task to db
	filenames := utils.GetFileName(param.Filename)
	if len(param.Filename) != 0 {
		_, outputFileName, filePathToSave, fileNameToSave, fileSizeToSave, outputFileSize = extractFilesToSave(&entity.TaskEntity{}, param.Filename, filenames, param.Filesize)
	}

	err = taskDAO.UpdateDoneTask(ctx, param.TaskID, outputFileName, outputFileSize, param.Filename, filenames, param.Filesize)
	if err != nil {
		return model.UpdateTaskSuccessResult{}, err
	}

	// STEP 2: Initiate new workflow execution if can
	childtask, err := taskDAO.GetReadyChildrenTaskByTaskID(ctx, param.TaskID)
	if err != nil {
		return model.UpdateTaskSuccessResult{}, err
	}

	wo := client.StartWorkflowOptions{
		TaskQueue: "task-queue-name",
	}
	for i := 0; i < len(childtask); i++ {
		res, err := u.tempCli.ExecuteWorkflow(ctx, wo, model.ExecuteTaskWfName, ExecuteTaskParam{Task: childtask[i]})
		if err != nil {
			//log.Err(err).Msg("[Wel logic internal] Unable to call GrantRoleWorkflow")
			return model.UpdateTaskSuccessResult{}, err
		}
		if err := res.Get(ctx, nil); err != nil {
			//log.Err(err).Msg("[Wel logic internal] GrantRoleWorkflow failed")
		}

		//log.Info().Str("Workflow", we.GetID()).Str("runID=", we.GetRunID()).Msg("dispatched")
	}

	// STEP 3: Return task's files
	task, err := taskDAO.GetTaskByTaskID(ctx, param.TaskID)
	if err != nil {
		return model.UpdateTaskSuccessResult{}, err
	}

	res = model.UpdateTaskSuccessResult{
		UserName:    task.UserName,
		RunUUID:     task.RunID.String(),
		ProjectID:   task.ProjectID.String(),
		ProjectPath: "/", // "default value"
		TaskID:      task.TaskID,
		TaskUUID:    task.ID.String(),
		TaskName:    task.Name,
		Path:        filePathToSave,
		Filename:    fileNameToSave,
		Filesize:    fileSizeToSave,
	}

	return res, nil
}

func (u *HeimdallTemporal) GetTaskByTaskIDAct(ctx context.Context, taskID string) (task entity.TaskEntity, err error) {
	taskDAO := repository.GetTaskDAO()
	return taskDAO.GetTaskByTaskID(ctx, taskID)
}

func (u *HeimdallTemporal) UpdateTaskFailAct(ctx context.Context, param model.UpdateTaskFailParam) (err error) {
	taskDAO := repository.GetTaskDAO()
	err = taskDAO.UpdateFailTask(ctx, param.TaskID)

	// try to delete redundance task
	return err
}

func extractFilesToSave(task *entity.TaskEntity, files, filenames []string, fileSize []int64) (outputFilePath, outputFileName, filePathToSave, fileNameToSave []string, fileSizeToSave, outputFileSize []int64) {
	if len(files) == 0 {
		return
	}
	fileMap := make(map[string]int64)
	for i := 0; i < len(files); i++ {
		for j := 0; j < len(task.OutputRegex); j++ {
			ok, _ := filepath.Match(task.OutputRegex[j], filenames[i])
			if ok {
				outputFilePath = append(outputFilePath, files[i])
				outputFileName = append(outputFileName, filenames[i])
				outputFileSize = append(outputFileSize, fileSize[i])
				filePathToSave = append(filePathToSave, files[i])
				fileNameToSave = append(fileNameToSave, filenames[i])
				fileSizeToSave = append(fileSizeToSave, fileSize[i])
				break
			}
		}
		fileMap[files[i]] = fileSize[i]
	}

	// check secondary file
	outputDir := filepath.Dir(files[0])
	secondaryFile := task.Output2ndFiles
	var counter []int
	var remainder []string

	for j1 := 0; j1 < len(secondaryFile); j1++ {
		count := strings.Count(secondaryFile[j1], "^")
		remain := strings.ReplaceAll(secondaryFile[j1], "^", "")
		counter = append(counter, count)
		remainder = append(remainder, remain)
	}

	for i := 0; i < len(filenames); i++ {
		fileElement := strings.Split(filenames[i], ".")
		for k := 0; k < len(counter); k++ {
			var secondaryFileFirstPath string
			for k1 := 0; k1 < len(fileElement)-counter[k]; k1++ {
				if k1 != len(fileElement)-counter[k]-1 {
					secondaryFileFirstPath = secondaryFileFirstPath + fileElement[k1] + "."
				} else {
					secondaryFileFirstPath = secondaryFileFirstPath + fileElement[k1]
				}
			}
			secondaryFileName := outputDir + "/" + secondaryFileFirstPath + remainder[k]

			if _, ok := fileMap[secondaryFileName]; ok {
				filePathToSave = append(filePathToSave, secondaryFileName)
				fileSizeToSave = append(fileSizeToSave, fileMap[secondaryFileName])
			}
		}
	}
	return
}

func transformToTaskDTO(task entity.TaskEntity) (res model.TaskDTO, err error) {
	var command string = ""
	for i := 0; i < len(task.Command); i++ {
		if i != len(task.Command)-1 {
			command = command + task.Command[i] + " "
		} else {
			command = command + task.Command[i]
		}
	}
	var param []*model.ParamWithRegex
	err = json.Unmarshal(task.ParamsWithRegex, &param)
	if err != nil {
		return model.TaskDTO{}, err
	}

	return model.TaskDTO{
		TaskID:          task.TaskID,
		TaskUUID:        task.ID.String(),
		Command:         command,
		ParamsWithRegex: param,
		OutputRegex:     task.OutputRegex,
		Output2ndFiles:  task.Output2ndFiles,
		DockerImage:     task.DockerImage,
	}, nil
}
