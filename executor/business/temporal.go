package business

import (
	"context"
	"time"

	"workflow/executor/core"
	"workflow/workflow-utils/model"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const ()

type ExecutorTemporal struct {
	tempCli   client.Client
	wfWorker  worker.Worker
	actWorker worker.Worker
	lg        *core.LogFormat
}

var executorTemp = &ExecutorTemporal{}

func RunTemporalDaemon(parentCtx context.Context) (fn model.Daemon, err error) {
	lg := core.GetLogger()
	lg.Info("Starting Temporal daemon")

	mainConf := core.GetMainConfig()
	sleepContext(parentCtx, time.Duration(mainConf.FailOverTime)*time.Second)
	if !core.GetSyncFlag() {
		core.SetSyncFlag()
	}

	if executorTemp.tempCli == nil {
		c, err := client.NewClient(client.Options{})
		if err != nil {
			lg.Fatalf("unable to create Temporal client", err)
		}
		executorTemp.tempCli = c
	}
	executorTemp.lg = lg
	err = executorTemp.RegisterWorker()
	if err != nil {
		lg.Fatalf("unable to create Temporal client", err)
	}

	fn = func() {
		<-parentCtx.Done()
		executorTemp.wfWorker.Stop()
		executorTemp.actWorker.Stop()
		executorTemp.tempCli.Close()

		lg.Info("Shutting down Temporal daemon")
	}

	return fn, nil
}

func GetExecutorTemporal() *ExecutorTemporal {
	return executorTemp
}

// Service implementation
func SetExecutorTemporal(cli client.Client) {
	executorTemp = &ExecutorTemporal{
		tempCli: cli,
	}
}

func (e *ExecutorTemporal) RegisterWorker() (err error) {
	// TODO: fix this after you have config
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize: 1,
	}

	// TODO: add task queue name
	e.wfWorker = worker.New(e.tempCli, model.BifrostExWf, workerOptions)
	// register workflow
	e.wfWorker.RegisterWorkflowWithOptions(e.DoneTaskWf, workflow.RegisterOptions{Name: model.DoneTaskWfName})
	e.wfWorker.RegisterWorkflowWithOptions(e.FailTaskWf, workflow.RegisterOptions{Name: model.FailTaskWfName})

	e.actWorker = worker.New(e.tempCli, model.BifrostExAct, workerOptions)
	// register activity
	e.actWorker.RegisterActivityWithOptions(e.ExecuteTaskAct, activity.RegisterOptions{Name: model.ExecuteTaskActName})
	e.actWorker.RegisterActivityWithOptions(e.DeleteTaskAct, activity.RegisterOptions{Name: model.DeleteTaskActName})

	// TODO: add LOGGGG
	if err := e.wfWorker.Start(); err != nil {
		e.lg.Error(err.Error())
		return err
	}

	if err := e.actWorker.Start(); err != nil {
		e.lg.Error(err.Error())
		return err
	}

	return nil
}

func (e *ExecutorTemporal) DoneTaskWf(ctx workflow.Context, param model.UpdateTaskSuccessParam) (err error) {
	e.lg.Info("[DoneTaskWf] Start Done Task workflow for " + param.TaskID)
	core.DecreaseJobCount()

	// STEP 1: Update task success to heimdall
	var res = model.UpdateTaskSuccessResult{}
	e.lg.Info("[DoneTaskWf] Execute update task success activity for " + param.TaskID)
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    500,
	}
	options := workflow.ActivityOptions{
		TaskQueue:           model.BifrostHeimAct,
		StartToCloseTimeout: 2 * time.Second,
		RetryPolicy:         retrypolicy,
	}
	ctx1 := workflow.WithActivityOptions(ctx, options)

	future := workflow.ExecuteActivity(ctx1, model.UpdateTaskSuccessActName, param)
	if err = future.Get(ctx, &res); err != nil {
		e.lg.Error(err.Error())
		return
	}

	// STEP 2: Push files that needed to be saved to valkyrire
	retrypolicy = &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    500,
	}
	options = workflow.ActivityOptions{
		TaskQueue:           model.BifrostValAct,
		StartToCloseTimeout: 2 * time.Second,
		RetryPolicy:         retrypolicy,
	}
	ctx2 := workflow.WithActivityOptions(ctx, options)
	e.lg.Info("[DoneTaskWf] Execute save generated file activity for " + param.TaskID)
	future = workflow.ExecuteActivity(ctx2, model.SaveGeneratedFileActName, res)
	if err = future.Get(ctx, nil); err != nil {
		e.lg.Error(err.Error())
		return
	}

	e.lg.Info("[DoneTaskWf] End Done Task workflow for " + param.TaskID)
	return nil
}

func (e *ExecutorTemporal) FailTaskWf(ctx workflow.Context, param model.UpdateTaskSuccessParam) (err error) {
	e.lg.Info("[FailTasktWf] Start Fail Task workflow for " + param.TaskID)
	core.DecreaseJobCount()
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    500,
	}
	options := workflow.ActivityOptions{
		TaskQueue:           model.BifrostHeimAct,
		StartToCloseTimeout: 2 * time.Second,
		RetryPolicy:         retrypolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	e.lg.Info("[FailTasktWf] Start update task failed activity for " + param.TaskID)
	future := workflow.ExecuteActivity(ctx, model.UpdateTaskFailActName, param)
	if err = future.Get(ctx, nil); err != nil {
		e.lg.Error(err.Error())
		return
	}

	e.lg.Info("[FailTasktWf] End Fail Task workflow for " + param.TaskID)
	return err
}

// TODO: add workflow def name and run id
func (e *ExecutorTemporal) ExecuteTaskAct(ctx context.Context, param model.ExecuteTaskParam) (res model.ExecuteTaskResult, err error) {
	// check for tasks threshold
	mainConf := core.GetMainConfig()
	for {
		if core.IsGoodToGo(mainConf.MaximumConcurrentJob) {
			break
		}
	}

	// add log here
	err = CreateK8SJob(ctx, &param.Task, e.lg, "", "")
	if err != nil {
		e.lg.Error(err.Error())
		core.DecreaseJobCount()
		return model.ExecuteTaskResult{
			TimeStamp: time.Now(),
			Created:   false,
		}, nil
	}
	res = model.ExecuteTaskResult{
		TimeStamp: time.Now(),
		Created:   true,
	}
	return res, nil
}

func (e *ExecutorTemporal) DeleteTaskAct(ctx context.Context, taskID string) (err error) {
	err = core.DeleteK8SJob(context.Background(), taskID, true)
	return err
}
