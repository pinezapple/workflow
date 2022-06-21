package business

import (
	"context"
	"time"

	"github.com/google/uuid"
	"workflow/valkyrie/core"
	"workflow/valkyrie/dao"
	valModel "workflow/valkyrie/model"
	"workflow/workflow-utils/model"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const ()

type ValkyrieTemporal struct {
	tempCli client.Client
	worker  worker.Worker
	//lg      *core.LogFormat
}

// Service implementation
func CreateValkyrieTemporal(cli client.Client) *ValkyrieTemporal {
	return &ValkyrieTemporal{
		tempCli: cli,
	}
}

func (e *ValkyrieTemporal) RegisterWorker() (err error) {
	// TODO: fix this after you have config
	workerOptions := worker.Options{
		MaxConcurrentWorkflowTaskExecutionSize: 1000,
	}
	// TODO: add task queue name
	e.worker = worker.New(e.tempCli, "your_task_queue_name", workerOptions)

	// register activity
	e.worker.RegisterActivityWithOptions(e.SaveGeneratedFileAct, activity.RegisterOptions{Name: model.SaveGeneratedFileActName})

	// TODO: add LOGGGG
	if err := e.worker.Start(); err != nil {
		return err
	}
	return nil
}

func (v *ValkyrieTemporal) SaveGeneratedFileAct(ctx context.Context, req model.UpdateTaskSuccessResult) (err error) {
	db := core.GetDBObj()
	mainConf := core.GetMainConfig()
	gfDAO := dao.GetGeneratedFileDAO()

	for i := 0; i < len(req.Filename); i++ {
		f := &valModel.GeneratedFile{
			FileUUID:      uuid.New().String(),
			UserID:        req.UserName,
			RunID:         req.RunID,
			RunUUID:       req.RunUUID,
			RunName:       req.RunName,
			ProjectID:     req.ProjectID,
			ProjectPath:   req.ProjectPath,
			TaskID:        req.TaskID,
			TaskUUID:      req.TaskUUID,
			TaskName:      req.TaskName,
			Path:          req.Path[i],
			Filename:      req.Filename[i],
			Filesize:      req.Filesize[i],
			UploadSuccess: false,
			DoneRun:       false,
			CreatedAt:     time.Now(),
		}

		err := gfDAO.SaveFile(ctx, db, f)
		if err != nil {
			logger.Errorf("Save generated file info to db err : %s", err.Error())
			return err
		}
	}

	if req.LastTask {
		expiredCloudTime := time.Now().Add(time.Duration(mainConf.NormalFileTTL) * time.Hour)

		err := gfDAO.UpdateCloudExpiredTime(ctx, db, req.UserName, req.RunUUID, expiredCloudTime)
		if err != nil {
			logger.Error(err.Error())
		}

		err = gfDAO.UpdateDoneRun(ctx, db, req.UserName, req.RunUUID)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	return nil
}
