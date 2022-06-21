// Package services call Scheduler Service via Kafka
package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"workflow/heimdall/core"
	"workflow/heimdall/services/dto"
	"workflow/heimdall/utils"
	utilmodel "workflow/workflow-utils/model"
)

// -----------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- SEND RUN TO SCHEDULER ----------------------------------------------------------

var (
	userRunExceededMsg   = "user has already submitted a DAG"
	ErrUserRunExceeded   = errors.New("the user's run number is exceeded")
	totalTaskExceededMsg = "too many tasks are in the system"
	ErrTotalTaskExceeded = errors.New("the task number is exceeded the threshold")
)

type CreateRunWorkflowForm struct {
	UserID         int                    `json:"user_id"`
	Retry          int                    `json:"retry"`
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Content        string                 `json:"content"`
	Steps          []*dto.WorkflowStep    `json:"steps"`
	WorkflowParams map[string]interface{} `json:"workflow_params"`
}

type schedulerService struct {
}

var (
	schedulerSrv = new(schedulerService)
)

func GetSchedulerService() *schedulerService {
	return schedulerSrv
}

func (scheSrv *schedulerService) GetAddress() (addr string) {
	scheConfig := core.GetConfig().Scheduler
	return scheConfig.Host + ":" + scheConfig.Port
}

func (scheSrv *schedulerService) SendRun(c context.Context, run *dto.Run, timeout int64) (err error) {
	schedRun := &dto.SchedRun{
		WorkflowID:  run.WorkflowID.String(),
		RunID:       run.RunID,
		ProjectID:   run.ProjectID.String(),
		ProjectPath: run.ProjectPath,
		RunName:     run.RunName,
		UserName:    run.UserName,
	}

	for i := range run.Tasks {
		schedTask := &dto.SchedTask{
			ID:              run.Tasks[i].ID.String(),
			TaskID:          run.Tasks[i].TaskID,
			RunID:           run.RunID,
			ProjectID:       run.Tasks[i].ProjectID.String(),
			ProjectPath:     run.Tasks[i].ProjectPath,
			RunUUID:         run.Tasks[i].RunID.String(),
			TaskName:        run.Tasks[i].TaskName,
			IsBoundary:      run.Tasks[i].IsBoundary,
			StepName:        run.Tasks[i].StepName,
			UserName:        run.Tasks[i].UserName,
			Command:         run.Tasks[i].Command,
			ParamsWithRegex: run.Tasks[i].ParamsWithRegex,
			OutputRegex:     run.Tasks[i].OutputRegex,
			Output2ndFiles:  run.Tasks[i].Output2ndFiles,
			ParentTasksID:   run.Tasks[i].ParentTasksID,
			ChildrenTasksID: run.Tasks[i].ChildrenTasksID,
			OutputLocation:  run.Tasks[i].OutputLocation,
			DockerImage:     run.Tasks[i].DockerImage,
			QueueLevel:      run.Tasks[i].QueueLevel,
		}
		schedRun.Tasks = append(schedRun.Tasks, schedTask)
	}

	sendData, err := json.Marshal(schedRun)
	if err != nil {
		return err
	}

	// printSchedRun, _ := json.MarshalIndent(schedRun, "", "    ")
	// fmt.Printf("\t\t SchedRun \n%s\n", string(printSchedRun))

	resp, status, err := utils.PostJSON(c, timeout, scheSrv.GetAddress()+"/receiveTasks", sendData)
	if err != nil {
		return err
	}

	if status != 200 {
		var respData = new(utilmodel.Response)
		err := json.Unmarshal(resp, respData)
		if err != nil {
			logger.Errorf("unmarshal scheduler response error: %v", err)
			return fmt.Errorf("Error with status %d. Message %s", status, string(resp))
		}
		switch {
		case strings.Contains(respData.Error.Message, userRunExceededMsg):
			return ErrUserRunExceeded

		case strings.Contains(respData.Error.Message, totalTaskExceededMsg):
			return ErrTotalTaskExceeded
		}
		return fmt.Errorf("schedule service error: %s", string(resp))
	}

	//NOTE: Remove because scheduler send update req before response.
	// if err = GetTaskStatusUpdater().UpdateStatusSent(c, strconv.Itoa(run.RunID)); err != nil {
	// 	return err
	// }

	return nil
}
