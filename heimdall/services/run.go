package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vfluxus/heimdall/core"
	"github.com/vfluxus/heimdall/repository"
	"github.com/vfluxus/heimdall/repository/entity"
	"github.com/vfluxus/heimdall/services/dto"
	"github.com/vfluxus/heimdall/webserver/forms"
)

// RunService interface
type RunService interface {
	GetRuns(ctx *gin.Context, pageSize int, pageToken int, filter map[string][]string) ([]forms.RunStatusDto, int64, error)
	CreateRun(ctx *gin.Context, runForm *forms.WorkflowRunForm) (forms.RunDto, error)
	GetRun(ctx *gin.Context, id uuid.UUID) (forms.RunDto, error)
	GetRunStatus(ctx *gin.Context, id uuid.UUID) (forms.RunStatusDto, error)
}

// GetRunService return Run service implement
func GetRunService() RunService {
	return runServiceImpl{}
}

type runServiceImpl struct{}

func (service runServiceImpl) GetRuns(ctx *gin.Context, pageSize, pageToken int, filter map[string][]string) (runs []forms.RunStatusDto, total int64, err error) {
	listRuns, total, err := repository.GetDAO().GetRuns(ctx, pageSize, pageToken, filter)
	if err != nil {
		return
	}
	for _, runEntity := range listRuns {
		var run = convertRunEntity2RunStatusDto(runEntity)
		runs = append(runs, run)
	}
	return
}

func (service runServiceImpl) CreateRunWorkflow(ctx *gin.Context, runForm *forms.WorkflowRunForm) (runDto forms.RunDto, err error) {
	dbDAO := repository.GetDAO()
	// FIXME temporary use WorkflowURL as Workflow UUID to select from db
	id, err := uuid.Parse(runForm.WorkflowURL)
	if err != nil {
		logger.Errorf("Convert workflow id failed: %s", runForm.WorkflowURL)
		return
	}

	workflow, err := dbDAO.GetWorkflow(ctx, id)
	if err != nil {
		logger.Errorf("Retrive workflow from db error: %s", err.Error())
		return
	}

	request := make(map[string]interface{})
	request["workflow_params"] = runForm.WorkflowParams
	request["workflow_type"] = runForm.WorkflowType
	request["workflow_type_version"] = runForm.WorkflowTypeVersion
	request["tags"] = runForm.Tags
	request["workflow_engine_parameters"] = runForm.WorkflowEngineParameters
	request["workflow_url"] = runForm.WorkflowURL
	request["workflow_attachments"] = runForm.WorkflowAttachments

	runRequest, err := json.Marshal(request)
	if err != nil {
		return
	}
	runEntity := &entity.RunEntity{
		ID: uuid.New(),
		// Description string
		// Tags       []byte `gorm:"type:jsonb"`
		Request:     runRequest,
		WorkflowID:  workflow.ID,
		UserName:    fmt.Sprintf("%v", ctx.Value("UserName")),
		State:       core.StateUnknown,
		ProjectID:   workflow.ProjectID,
		ProjectPath: "/", // TODO: HARD CODE FIX
	}

	// A run is should be created before sending it to the scheduler
	// because the scheduler will update the task status right after
	// it receives a new run.
	if err = dbDAO.CreateRun(ctx, runEntity); err != nil {
		return
	}

	var steps []*dto.WorkflowStep
	for _, step := range workflow.Steps {
		workflowStep := &dto.WorkflowStep{
			Name:    step.Name,
			Content: step.Content,
		}
		steps = append(steps, workflowStep)
	}

	tfReq := dto.TransformRequest{
		RunIndex: runEntity.RunIndex,
		UserName: fmt.Sprintf("%v", ctx.Value("UserName")),
		Name:     workflow.Name,
		Content:  workflow.Content,
		Params:   runForm.WorkflowParams,
		Steps:    steps,
	}

	//logger.Infof("Transform request: %v", tfReq)

	tfRes, err := GetTransformerService().Transform(ctx, tfReq)
	if err != nil {
		_ = dbDAO.DeleteRun(ctx, runEntity)
		return
	}

	// updateRunEntity add tasks into the runentity
	if err = updateRunEntity(ctx, runEntity, tfRes, false); err != nil {
		return
	}

	runDto, err = convertRunEntity2Dto(runEntity)
	if err != nil {
		logger.Errorf("Convert to run dto error: %s", err.Error())
		return
	}

	// add run metadata
	tfRes.Data.WorkflowID = workflow.ID
	tfRes.Data.ProjectID = runEntity.ProjectID
	tfRes.Data.ProjectPath = runEntity.ProjectPath

	// update the run with the transformed data
	err = dbDAO.UpdateRun(ctx, runEntity)
	if err != nil {
		return
	}

	// TODO: start to init workflow starting jobs here

	return runDto, nil
}

func (service runServiceImpl) CreateRun(ctx *gin.Context, runForm *forms.WorkflowRunForm) (runDto forms.RunDto, err error) {
	dbDAO := repository.GetDAO()
	// FIXME temporary use WorkflowURL as Workflow UUID to select from db
	id, err := uuid.Parse(runForm.WorkflowURL)
	if err != nil {
		logger.Errorf("Convert workflow id failed: %s", runForm.WorkflowURL)
		return
	}

	workflow, err := dbDAO.GetWorkflow(ctx, id)
	if err != nil {
		logger.Errorf("Retrive workflow from db error: %s", err.Error())
		return
	}

	request := make(map[string]interface{})
	request["workflow_params"] = runForm.WorkflowParams
	request["workflow_type"] = runForm.WorkflowType
	request["workflow_type_version"] = runForm.WorkflowTypeVersion
	request["tags"] = runForm.Tags
	request["workflow_engine_parameters"] = runForm.WorkflowEngineParameters
	request["workflow_url"] = runForm.WorkflowURL
	request["workflow_attachments"] = runForm.WorkflowAttachments

	runRequest, err := json.Marshal(request)
	if err != nil {
		return
	}
	runEntity := &entity.RunEntity{
		ID: uuid.New(),
		// Description string
		// Tags       []byte `gorm:"type:jsonb"`
		Request:     runRequest,
		WorkflowID:  workflow.ID,
		UserName:    fmt.Sprintf("%v", ctx.Value("UserName")),
		State:       core.StateUnknown,
		ProjectID:   workflow.ProjectID,
		ProjectPath: "/", // TODO: HARD CODE FIX
	}

	// A run is should be created before sending it to the scheduler
	// because the scheduler will update the task status right after
	// it receives a new run.
	if err = dbDAO.CreateRun(ctx, runEntity); err != nil {
		return
	}

	var steps []*dto.WorkflowStep
	for _, step := range workflow.Steps {
		workflowStep := &dto.WorkflowStep{
			Name:    step.Name,
			Content: step.Content,
		}
		steps = append(steps, workflowStep)
	}

	tfReq := dto.TransformRequest{
		RunIndex: runEntity.RunIndex,
		UserName: fmt.Sprintf("%v", ctx.Value("UserName")),
		Name:     workflow.Name,
		Content:  workflow.Content,
		Params:   runForm.WorkflowParams,
		Steps:    steps,
	}

	//logger.Infof("Transform request: %v", tfReq)

	tfRes, err := GetTransformerService().Transform(ctx, tfReq)
	if err != nil {
		_ = dbDAO.DeleteRun(ctx, runEntity)
		return
	}
	// TODO OOOOOOOO: create task from here

	// updateRunEntity add tasks into the runentity
	if err = updateRunEntity(ctx, runEntity, tfRes, false); err != nil {
		return
	}

	runDto, err = convertRunEntity2Dto(runEntity)
	if err != nil {
		logger.Errorf("Convert to run dto error: %s", err.Error())
		return
	}

	if err != nil {
		logger.Errorf("Get workflow %s error: %s", runEntity.WorkflowID, err.Error())
		return
	}

	// add run metadata
	tfRes.Data.WorkflowID = workflow.ID
	tfRes.Data.ProjectID = runEntity.ProjectID
	tfRes.Data.ProjectPath = runEntity.ProjectPath

	// add task metadata
	var idUUID = make(map[string]uuid.UUID)
	for taskIndex := range runEntity.Tasks {
		idUUID[runEntity.Tasks[taskIndex].TaskID] = runEntity.Tasks[taskIndex].ID
	}
	for taskIndex := range tfRes.Data.Tasks {
		id, ok := idUUID[tfRes.Data.Tasks[taskIndex].TaskID]
		if !ok {
			return runDto, fmt.Errorf("Can not find task ID in runEntity: %s", tfRes.Data.Tasks[taskIndex].TaskID)
		}
		tfRes.Data.Tasks[taskIndex].ID = id
		tfRes.Data.Tasks[taskIndex].RunID = runEntity.ID
		tfRes.Data.Tasks[taskIndex].ProjectID = runEntity.ProjectID
		tfRes.Data.Tasks[taskIndex].ProjectPath = runEntity.ProjectPath
	}

	// update the run with the transformed data
	err = dbDAO.UpdateRun(ctx, runEntity)
	if err != nil {
		return
	}

	err = GetSchedulerService().SendRun(ctx, &tfRes.Data, 100)
	if err != nil {
		// if the scheduler doesn't accept the run, delete it
		if err := dbDAO.DeleteRun(ctx, runEntity); err != nil {
			logger.Errorf("CreateRun - delete the run: %s", runEntity.ID)
			logger.Errorf("Delete run error: %s", err.Error())
		}
		return runDto, err
	}

	return runDto, nil
}

func (service runServiceImpl) GetRun(ctx *gin.Context, id uuid.UUID) (run forms.RunDto, err error) {
	runEntity, err := repository.GetDAO().GetRun(ctx, id)
	if err != nil {
		return
	}

	// remove logical node
	var t = make([]entity.TaskEntity, 0, len(runEntity.Tasks))
	for i := range runEntity.Tasks {
		if runEntity.Tasks[i].IsBoundary {
			continue
		}

		t = append(t, runEntity.Tasks[i])
	}
	runEntity.Tasks = t
	run, err = convertRunEntity2Dto(&runEntity)
	return
}

func (service runServiceImpl) GetRunStatus(ctx *gin.Context, id uuid.UUID) (runStatus forms.RunStatusDto, err error) {
	runEntity, err := repository.GetDAO().GetRun(ctx, id)
	if err != nil {
		return
	}
	runStatus = convertRunEntity2RunStatusDto(runEntity)
	return
}

func updateRunEntity(ctx *gin.Context, runEntity *entity.RunEntity, tfRes dto.TransformResponse, generateID bool) (err error) {
	// ---------------- Create List Task Entity ----------------------------------
	var listTasks []entity.TaskEntity
	for _, task := range tfRes.Data.Tasks {
		commands := strings.Split(strings.Trim(task.Command, " "), " ")
		cmds := make([]string, len(commands))
		copy(cmds, commands)
		// marshal to save
		paramsWithRegexByte, err := json.Marshal(task.ParamsWithRegex)
		if err != nil {
			return err
		}

		taskEntity := entity.TaskEntity{
			ID:          uuid.New(),
			TaskID:      task.TaskID,
			RunID:       runEntity.ID,
			RunIndex:    runEntity.RunIndex,
			ProjectID:   runEntity.ProjectID,
			ProjectPath: runEntity.ProjectPath,
			// Name:
			// Description
			IsBoundary: task.IsBoundary,
			StepName:   task.StepName,
			UserName:   fmt.Sprintf("%v", ctx.Value("UserName")),
			Command:    cmds,
			// Inputs
			// Outputs
			// Resource:  "",
			// Executors: "",
			// Logs
			OutputRegex:     task.OutputRegex,
			DockerImage:     task.DockerImage,
			Output2ndFiles:  task.Output2ndFiles,
			ParamsWithRegex: paramsWithRegexByte,
			ParentTasksID:   task.ParentTasksID,
			ChildrenTasksID: task.ChildrenTasksID,
			OutputLocation:  task.OutputLocation,
			State:           core.StateUnknown,
			// State:           task.State,
			// StartedTime
			// EndTime
		}

		if generateID {
			taskEntity.ID = uuid.New()
		}
		listTasks = append(listTasks, taskEntity)
	}

	runEntity.Tasks = listTasks
	return
}

func convertRunEntity2RunStatusDto(runEntity entity.RunEntity) forms.RunStatusDto {
	var runRequest = forms.RunRequestDto{}
	if err := json.Unmarshal(runEntity.Request, &runRequest); err != nil {
		logger.Errorf("Conver run request error: %s", err.Error())
		return forms.RunStatusDto{}
	}

	var runStatus = forms.RunStatusDto{
		ID:        runEntity.ID,
		State:     runEntity.State,
		User:      runEntity.UserName,
		Request:   runRequest,
		StartTime: runEntity.CreatedAt,
		EndTime:   runEntity.EndTime.Time,
	}
	return runStatus
}

func convertRunEntity2Dto(runEntity *entity.RunEntity) (forms.RunDto, error) {
	var runRequest = forms.RunRequestDto{}
	if err := json.Unmarshal(runEntity.Request, &runRequest); err != nil {
		logger.Errorf("Conver run request error: %s", err.Error())
		return forms.RunDto{}, err
	}

	runDto := forms.RunDto{
		ID:      runEntity.ID,
		Request: runRequest,
		State:   runEntity.State,
	}

	if runEntity.RunLog != nil {
		var runLog = forms.RunLogDto{}
		if err := json.Unmarshal(runEntity.RunLog, &runLog); err != nil {
			logger.Errorf("Convert run log error: %s", err.Error())
			return forms.RunDto{}, err
		}

		runDto.Log = runLog
	}

	var tasks []forms.SimpleTaskDto
	var outputs map[string]interface{} = make(map[string]interface{})
	for _, task := range runEntity.Tasks {
		var taskName = strings.Split(task.TaskID, "-")[len(strings.Split(task.TaskID, "-"))-1]
		simpleTask := forms.SimpleTaskDto{
			ID:        task.ID,
			State:     task.State,
			Name:      taskName, // TODO(tuandn8) change to use task name
			StartedAt: task.StartedTime,
			EndAt:     task.EndTime,
		}
		tasks = append(tasks, simpleTask)

		// if len(task.Outputs) > 0 {
		// 	if err = json.Unmarshal(task.Outputs, &outputs); err != nil {
		// 		return
		// 	}
		// } else {
		// TODO(tuandn8) Temporary copy data from OutputLocation to Outputs
		var taskOutputs []forms.TaskOutputDto
		for _, file := range task.OutputLocation {
			var taskOutput = forms.TaskOutputDto{
				Name: file,
				// Description :
				// URL: file,
				// Path        :
				Type: "FILE",
			}
			taskOutputs = append(taskOutputs, taskOutput)
		}
		outputs[taskName] = struct {
			Outputs []forms.TaskOutputDto `json:"outputs"`
		}{
			Outputs: taskOutputs,
		}
	}
	runDto.Tasks = tasks
	runDto.Outputs = outputs

	return runDto, nil
}
