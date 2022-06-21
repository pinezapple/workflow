package business

import (
	"context"
	"time"

	"github.com/vfluxus/workflow-utils/model"
	"github.com/vfluxus/workflow/executor/core"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const ()

type ExecutorTemporal struct {
	tempCli client.Client
	worker  worker.Worker
	lg      *core.LogFormat
}

// Service implementation
func CreateExecutorTemporal(cli client.Client) *ExecutorTemporal {
	return &ExecutorTemporal{
		tempCli: cli,
	}
}

func (e *ExecutorTemporal) RegisterWorker() (err error) {
	// TODO: fix this after you have config
	workerOptions := worker.Options{
		MaxConcurrentWorkflowTaskExecutionSize: 1000,
	}
	// TODO: add task queue name
	e.worker = worker.New(e.tempCli, "your_task_queue_name", workerOptions)

	// register workflow
	e.worker.RegisterWorkflowWithOptions(e.DoneTasktWf, workflow.RegisterOptions{Name: model.DoneTasktWfName})
	e.worker.RegisterWorkflowWithOptions(e.FailTasktWf, workflow.RegisterOptions{Name: model.FailTasktWfName})

	// register activity
	e.worker.RegisterActivityWithOptions(e.ExecuteTaskAct, activity.RegisterOptions{Name: model.ExecuteTaskActName})
	e.worker.RegisterActivityWithOptions(e.DeleteTaskAct, activity.RegisterOptions{Name: model.DeleteTaskActName})

	// TODO: add LOGGGG
	if err := e.worker.Start(); err != nil {
		e.lg.Error(err.Error())
		return err
	}
	return nil
}

func (e *ExecutorTemporal) DoneTasktWf(ctx workflow.Context, param model.UpdateTaskSuccessParam) (err error) {
	// STEP 1: Update task success to heimdall
	e.lg.Info("[DoneTaskWf] Start Done Task workflow")

	var res = model.UpdateTaskSuccessResult{}
	e.lg.Info("[DoneTaskWf] Execute update task success activity")
	future := workflow.ExecuteActivity(ctx, model.UpdateTaskSuccessActName, param)
	if err = future.Get(ctx, &res); err != nil {
		e.lg.Error(err.Error())
		return
	}

	// STEP 2: Push files that needed to be saved to valkyrire
	e.lg.Info("[DoneTaskWf] Execute save generated file activity")
	future = workflow.ExecuteActivity(ctx, model.SaveGeneratedFileActName, res)
	if err = future.Get(ctx, nil); err != nil {
		e.lg.Error(err.Error())
		return
	}

	e.lg.Info("[DoneTaskWf] End Done Task workflow")
	return nil
}

func (e *ExecutorTemporal) FailTasktWf(ctx workflow.Context, param model.UpdateTaskSuccessParam) (err error) {
	future := workflow.ExecuteActivity(ctx, model.UpdateTaskFailActName, param)
	if err = future.Get(ctx, nil); err != nil {
		e.lg.Error(err.Error())
		return
	}

	return err
}

// TODO: add workflow def name and run id
func (e *ExecutorTemporal) ExecuteTaskAct(ctx context.Context, param model.ExecuteTaskParam) (res model.ExecuteTaskResult, err error) {
	// check for tasks threshold
	// add log here
	err = CreateK8SJob(ctx, &param.Task, e.lg, "", "")
	if err != nil {
		return model.ExecuteTaskResult{}, err
	}
	res = model.ExecuteTaskResult{
		TimeStamp: time.Now(),
	}
	return res, nil
}

func (e *ExecutorTemporal) DeleteTaskAct(ctx context.Context, taskID string) (err error) {
	err, _ = core.DeleteK8SJob(context.Background(), taskID, true)
	return err
}
