package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vfluxus/heimdall/services"
	"github.com/vfluxus/heimdall/webserver/forms"
)

// HandleGETTasks ...
// @Summary Get tasks info
// @Description Get tasks info
// @Produce json
// @Param 	page_size 		query 	int 	true "page size"
// @Param 	page_token 		query 	int 	true "page token"
// @Param 	filter query 	query	string 	true "Filter. Split by ;"
// @Success 200 {object} forms.ListTasksDto "tasks"
// @Tags task
// @Router /tasks [GET]
func HandleGETTasks(c *gin.Context) {
	ok := AuthzRequest(c, "/workflow/run/task", "read", "heimdall")
	if ok == false {
		ResponseError(c, errors.New("You do not have permission"), http.StatusForbidden)
		return
	}

	pageSize, pageToken, filterMap, err := getFilterParam(c)
	if err != nil {
		return
	}

	tasks, total, err := services.GetTaskService().GetTasks(c, pageSize, pageToken, filterMap)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &forms.ListTasksDto{
		Tasks:         tasks,
		NextPageToken: strconv.Itoa(pageToken + 1),
		Total:         total,
	})

	return
}

// HandleGETTask ...
// @Summary Get task info by id
// @Description Get task info by id
// @Produce json
// @Param task_id path string  true "task id"
// @Success 200 {object} forms.TaskDto
// @Tags task
// @Router /tasks/:task_id [GET]
func HandleGETTask(c *gin.Context) {
	ok := AuthzRequest(c, "/workflow/run/task", "read", "heimdall")
	if ok == false {
		ResponseError(c, errors.New("You do not have permission"), http.StatusForbidden)
		return
	}

	task_id := c.Param("task_id")
	id, err := uuid.Parse(task_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", task_id)
		ResponseError(c, errors.New("Convert task id error"), http.StatusBadRequest)
		return
	}

	task, err := services.GetTaskService().GetTask(c, id)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, task)
	return
}

// HandlePOSTTask ...
// @Summary Create task
// @Description Create task
// @Produce json
// @Param createReq body forms.TaskFormDto true "task info"
// @Success 200 {string} id=string "uuid"
// @Tags task
// @Router /tasks [POST]
func HandlePOSTTask(c *gin.Context) {
	ok := AuthzRequest(c, "/workflow/run/task", "create", "heimdall")
	if ok == false {
		ResponseError(c, errors.New("You do not have permission"), http.StatusForbidden)
		return
	}

	var taskForm forms.TaskFormDto
	err := c.BindJSON(&taskForm)
	if err != nil {
		logger.Errorf("Error when parse taskdto: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	taskDto, err := services.GetTaskService().CreateTask(c, &taskForm)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": taskDto.ID,
	})
}
