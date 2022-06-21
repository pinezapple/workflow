package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vfluxus/heimdall/core"
	"github.com/vfluxus/heimdall/repository"
	"github.com/vfluxus/heimdall/repository/entity"
	"github.com/vfluxus/heimdall/webserver/forms"
)

// TaskService interface contains methods
type TaskService interface {
	GetTask(ctx *gin.Context, id uuid.UUID) (forms.TaskDto, error)
	GetTasks(ctx *gin.Context, pageSize int, pageToken int, filter map[string][]string) ([]forms.TaskDto, int64, error)
	CreateTask(ctx *gin.Context, taskForm *forms.TaskFormDto) (forms.TaskDto, error)
}

// GetTaskService returns implemented TaskService interface
func GetTaskService() TaskService {
	return taskServiceImpl{}
}

type taskServiceImpl struct{}

func (service taskServiceImpl) GetTasks(ctx *gin.Context, pageSize, pageToken int, filter map[string][]string) (tasks []forms.TaskDto, total int64, err error) {
	listTasks, total, err := repository.GetDAO().GetTasks(ctx, pageSize, pageToken, filter)
	if err != nil {
		return
	}

	for _, taskEntity := range listTasks {
		task, er := convertTaskEntity2Dto(ctx, &taskEntity)
		if er != nil {
			err = er
			logger.Errorf("%v", err.Error())
			return
		}

		tasks = append(tasks, task)
	}
	return
}

func (service taskServiceImpl) GetTask(ctx *gin.Context, id uuid.UUID) (task forms.TaskDto, err error) {
	taskEntity, err := repository.GetDAO().GetTask(ctx, id)
	if err != nil {
		return
	}

	task, err = convertTaskEntity2Dto(ctx, &taskEntity)
	return
}

func (service taskServiceImpl) CreateTask(ctx *gin.Context, taskForm *forms.TaskFormDto) (taskDto forms.TaskDto, err error) {
	taskEntity, err := convertTaskForm2Entity(ctx, taskForm, true)
	if err != nil {
		return forms.TaskDto{}, err
	}

	dbDAO := repository.GetDAO()
	id, err := uuid.Parse(taskForm.RunURL)
	if err != nil {
		logger.Errorf("Convert run id error: %s", taskForm.RunURL)
		return
	}
	run, err := dbDAO.GetRun(ctx, id)
	if err != nil {
		return forms.TaskDto{}, err
	}

	taskEntity.RunID = run.ID
	err = dbDAO.CreateTask(ctx, taskEntity)
	if err != nil {
		return forms.TaskDto{}, err
	}

	taskDto, err = convertTaskEntity2Dto(ctx, taskEntity)
	if err != nil {
		return forms.TaskDto{}, err
	}
	return taskDto, nil
}

func convertTaskForm2Entity(ctx *gin.Context, taskForm *forms.TaskFormDto, generateID bool) (*entity.TaskEntity, error) {
	taskEntity := &entity.TaskEntity{
		Name:        taskForm.Name,
		Description: taskForm.Description,
		UserName:    fmt.Sprintf("%v", ctx.Value("UserName")),
		Command:     taskForm.Executors[0].Command,
	}

	inputs, err := json.Marshal(taskForm.Inputs)
	if err != nil {
		return nil, err
	}

	outputs, err := json.Marshal(taskForm.Outputs)
	if err != nil {
		return nil, err
	}

	resources, err := json.Marshal(taskForm.Resources)
	if err != nil {
		return nil, err
	}

	executors, err := json.Marshal(taskForm.Executors)
	if err != nil {
		return nil, err
	}

	taskEntity.Inputs = inputs
	taskEntity.Outputs = outputs
	taskEntity.Resource = resources
	taskEntity.Executors = executors
	taskEntity.StartedTime = time.Now()

	if generateID {
		taskEntity.ID = uuid.New()
	}

	return taskEntity, nil
}

func convertTaskEntity2Dto(ctx *gin.Context, taskEntity *entity.TaskEntity) (taskDto forms.TaskDto, err error) {
	var inputs []forms.TaskInputDto
	var outputs []forms.TaskOutputDto
	var resources forms.TaskResourcesDto
	var executors []forms.TaskExecutorDto

	if len(taskEntity.Inputs) > 0 {
		if err = json.Unmarshal(taskEntity.Inputs, &inputs); err != nil {
			return
		}
	}

	if len(taskEntity.Outputs) > 0 {
		if err = json.Unmarshal(taskEntity.Outputs, &outputs); err != nil {
			return
		}
	} else {
		// TODO(tuandn8) Temporary copy data from OutputLocation to Outputs
		for _, file := range taskEntity.OutputLocation {
			var output = forms.TaskOutputDto{
				Name: file,
				// Description :
				URL: formDownloadURL(
					core.GetConfig().Valkyrie.Host,
					core.GetConfig().Valkyrie.Port,
					file,
					taskEntity.RunID.String(),
					taskEntity.ID.String()),
				// Path        :
				Type: "FILE",
			}
			outputs = append(outputs, output)
		}
	}

	if taskEntity.Resource != nil {
		if err = json.Unmarshal(taskEntity.Resource, &resources); err != nil {
			return
		}
	}

	if len(taskEntity.Executors) > 0 {
		if err = json.Unmarshal(taskEntity.Executors, &executors); err != nil {
			return
		}
	}

	var taskName string
	taskName = strings.Split(taskEntity.TaskID, "-")[len(strings.Split(taskEntity.TaskID, "-"))-1]

	taskDto = forms.TaskDto{
		ID:           taskEntity.ID,
		State:        taskEntity.State,
		Name:         taskName,
		Description:  taskEntity.Description,
		Inputs:       inputs,
		Outputs:      outputs,
		Resources:    resources,
		Executors:    executors,
		CreationTime: taskEntity.CreatedAt,
		StartedAt:    taskEntity.StartedTime,
		EndAt:        taskEntity.EndTime,
	}

	return
}

func formDownloadURL(valkyrieHost, valkyriePort, filename, runUUID, taskUUID string) string {
	if len(valkyriePort) != 0 {
		return fmt.Sprintf("%s:%s/files/%s?run_uuid=%s&task_uuid=%s", valkyrieHost, valkyriePort, filename, runUUID, taskUUID)
	}

	return fmt.Sprintf("%s/files/%s?run_uuid=%s&task_uuid=%s", valkyrieHost, filename, runUUID, taskUUID)
}
