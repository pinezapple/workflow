package services

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"workflow/heimdall/repository"
	"workflow/heimdall/repository/entity"
	"workflow/heimdall/webserver/forms"
)

// WorkflowService interface
type WorkflowService interface {
	// GetWorkflows return list workflows based on filter
	GetWorkflows(ctx *gin.Context, pageSize int, pageToken int,
		filter map[string][]string) ([]forms.WorkflowDto, int64, error)

	// GetWorkflowByID return workflow based on uuid
	GetWorkflowByID(ctx *gin.Context, id uuid.UUID) (forms.WorkflowDto, error)
	// CreateWorkflow creates a workflow and return newly created workflow
	CreateWorkflow(ctx *gin.Context, workflow forms.WorkflowForm) (forms.WorkflowDto, error)
	// UpdateWorkflow updates a workflow and return newly updated workflow
	UpdateWorkflow(ctx *gin.Context, id uuid.UUID, workflow forms.WorkflowForm) (forms.WorkflowDto, error)
	// DeleteWorkflow deletes a workflow based on uuid
	DeleteWorkflow(ctx *gin.Context, id uuid.UUID) error
	// GetWorkflowRuns return list runs of a workflow
	GetWorkflowRuns(ctx *gin.Context, username string, id uuid.UUID, pageSize int, pageToken int, filters map[string][]string) ([]forms.RunStatusDto, int64, error)
}

// GetWorkflowService return workflow service implement
func GetWorkflowService() WorkflowService {
	return workflowServiceImpl{}
}

type workflowServiceImpl struct{}

func (service workflowServiceImpl) GetWorkflows(ctx *gin.Context, pageSize int, pageToken int, filter map[string][]string) (dto []forms.WorkflowDto, total int64, err error) {
	listWorkflows, total, er := repository.GetDAO().GetWorkflows(ctx, pageSize, pageToken, filter)
	if er != nil {
		logger.Errorf("Error when get workflow dao: %s", er.Error())
		return
	}
	for _, workflowEntity := range listWorkflows {
		item := convertWorkflowEntity2Dto(&workflowEntity)
		dto = append(dto, item)
	}
	return
}

func (service workflowServiceImpl) GetWorkflowByID(ctx *gin.Context, id uuid.UUID) (dto forms.WorkflowDto, err error) {
	wfDB, err := repository.GetDAO().GetWorkflow(ctx, id)
	if err != nil {
		logger.Errorf("Error when get workflow by uuid: %s", err.Error())
		return
	}
	dto = convertWorkflowEntity2Dto(&wfDB)
	return
}

func (service workflowServiceImpl) GetWorkflowRuns(ctx *gin.Context, username string, id uuid.UUID, pageSize int, pageToken int,
	filters map[string][]string) (runs []forms.RunStatusDto, total int64, err error) {
	runsEntity, total, err := repository.GetDAO().GetRunsByWorkflowID(ctx, username, id, pageSize, pageToken, filters)
	if err != nil {
		logger.Errorf("Error when get workflow 's run: %s", err.Error())
		return
	}
	for _, runEntity := range runsEntity {
		var runStatusDto = convertRunEntity2RunStatusDto(runEntity)
		runs = append(runs, runStatusDto)
	}
	return
}

func (service workflowServiceImpl) CreateWorkflow(ctx *gin.Context,
	workflow forms.WorkflowForm) (forms.WorkflowDto, error) {

	workflowEntity, err := convertWorkflowForm2Entity(ctx, workflow, true)
	if err != nil {
		logger.Errorf("Convert workflow form to entity error: %s", err.Error())
		return forms.WorkflowDto{}, err
	}

	workflowEntity.ProjectID = workflow.ProjectID
	workflowEntity.ID = uuid.New()

	dbDAO := repository.GetDAO()
	err = dbDAO.CreateWorkflow(ctx, workflowEntity)
	if err != nil {
		return forms.WorkflowDto{}, err
	}

	return convertWorkflowEntity2Dto(workflowEntity), nil
}

func (service workflowServiceImpl) UpdateWorkflow(ctx *gin.Context,
	id uuid.UUID, workflow forms.WorkflowForm) (forms.WorkflowDto, error) {

	workflowEntity, err := convertWorkflowForm2Entity(ctx, workflow, false)
	if err != nil {
		logger.Errorf("Convert workflow form to entity error: %s", err.Error())
		return forms.WorkflowDto{}, err
	}
	workflowEntity.ID = id
	dbDAO := repository.GetDAO()
	err = dbDAO.UpdateWorkflow(ctx, id, workflowEntity)
	if err != nil {
		return forms.WorkflowDto{}, err
	}

	// query for return
	*workflowEntity, err = dbDAO.GetWorkflow(ctx, workflowEntity.ID)
	if err != nil {
		return forms.WorkflowDto{}, err
	}

	//logger.Infof("%v", *workflowEntity)

	return convertWorkflowEntity2Dto(workflowEntity), nil
}

func (service workflowServiceImpl) DeleteWorkflow(ctx *gin.Context, id uuid.UUID) error {

	dbDAO := repository.GetDAO()
	if err := dbDAO.DeleteWorkflow(ctx, id); err != nil {
		logger.Errorf("Error when delete workflow: %s", err.Error())
		return err
	}

	return nil
}

func convertWorkflowEntity2Dto(workflowEntity *entity.WorkflowEntity) (workflowDto forms.WorkflowDto) {
	workflowDto = forms.WorkflowDto{
		ID:          workflowEntity.ID,
		Name:        workflowEntity.Name,
		Class:       workflowEntity.Class,
		Description: workflowEntity.Description,
		Content:     workflowEntity.Content,
		Author:      workflowEntity.Author,
		CreatedAt:   workflowEntity.CreatedAt,
		UpdatedAt:   workflowEntity.UpdatedAt.Time,
		ProjectId:   workflowEntity.Project.ID,
		ProjectName: workflowEntity.Project.Name,
	}

	var tagsDto map[string]interface{}
	if err := json.Unmarshal(workflowEntity.Tags, &tagsDto); err != nil {
		logger.Errorf("Convert from []byte to map[string] error: %s", err.Error())
		return forms.WorkflowDto{}
	}
	workflowDto.Tags = tagsDto

	var listStepsDto []forms.WorkflowStepDto
	for _, step := range workflowEntity.Steps {
		var stepDto = forms.WorkflowStepDto{
			ID:        step.ID,
			Name:      step.Name,
			Content:   step.Content,
			CreatedAt: step.CreatedAt,
			UpdatedAt: step.UpdatedAt.Time,
		}
		listStepsDto = append(listStepsDto, stepDto)
	}
	workflowDto.Steps = listStepsDto
	return
}

func convertWorkflowForm2Entity(ctx *gin.Context, workflow forms.WorkflowForm,
	generateID bool) (workflowEntity *entity.WorkflowEntity, err error) {
	tags, err := json.Marshal(workflow.Tags)
	if err != nil {
		return nil, err
	}

	workflowEntity = &entity.WorkflowEntity{
		Name:        workflow.Name,
		Description: workflow.Description,
		Content:     workflow.Content,
		Class:       workflow.Class,
		Author:      "tungnt99",
		Tags:        tags,
	}
	if generateID {
		workflowEntity.ID = uuid.New()
	}

	var listStepEntity []entity.WorkflowStepEntity
	for _, step := range workflow.Steps {
		stepEntity := entity.WorkflowStepEntity{
			Name:    step.Name,
			Content: step.Content,
		}
		if generateID {
			stepEntity.ID = uuid.New()
		}

		listStepEntity = append(listStepEntity, stepEntity)
	}
	workflowEntity.Steps = listStepEntity
	return
}
