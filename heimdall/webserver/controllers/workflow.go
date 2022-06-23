package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"workflow/heimdall/core"
	"workflow/heimdall/services"
	"workflow/heimdall/webserver/forms"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var logger = core.GetLogger()

// HandleGETWorkflows handle get workflows
// @Summary Get workflows by param
// @Description Get workflows by param
// @Produce json
// @Param 	page_size 		query 	int 	true "page size"
// @Param 	page_token 		query 	int 	true "page token"
// @Param 	filter query 	query	string 	true "Filter. Split by ;"
// @Success 200 {object} forms.WorkflowsDto "Workflows information"
// @Tags workflow
// @Router /workflows [GET]
func HandleGETWorkflows(c *gin.Context) {
	pageSize, pageToken, filterMap, err := getFilterParam(c)
	if err != nil {
		return
	}

	workflowService := services.GetWorkflowService()
	workflows, total, err := workflowService.GetWorkflows(c, pageSize, pageToken, filterMap)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	response := &forms.WorkflowsDto{
		Workflows:     workflows,
		NextPageToken: strconv.Itoa(pageToken + 1),
		Total:         total,
	}
	c.JSON(http.StatusOK, response)
	return
}

// HandlePOSTWorkflow handle create workflow
// @Summary Create workflow
// @Description Create workflow
// @Produce json
// @Param createRequest body forms.WorkflowForm true "workflow info"
// @Success 200 {object} forms.WorkflowDto "Workflow info"
// @Tags workflow
// @Router /workflows [POST]
func HandlePOSTWorkflow(c *gin.Context) {
	var workflowForm forms.WorkflowForm
	err := c.Bind(&workflowForm)
	if err != nil {
		logger.Errorf("Parse workflow error: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	workflowService := services.GetWorkflowService()
	workflow, err := workflowService.CreateWorkflow(c, workflowForm)
	if err != nil {
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"workflow_id": workflow.ID})
	c.JSON(http.StatusOK, workflow)
	return
}

// HandleGETWorkflowByID handle get workflow by ID
// @Summary Get workflows by ID
// @Description Get workflows by ID
// @Produce json
// @Param workflow_id path int true "workflow id"
// @Success 200 {object} forms.WorkflowDto "workflow response"
// @Tags workflow
// @Router /workflows/:workflow_id [GET]
func HandleGETWorkflowByID(c *gin.Context) {
	workflow_id := c.Param("workflow_id")
	id, err := uuid.Parse(workflow_id)
	if err != nil {
		logger.Errorf("Convert workflow id failed: %s", workflow_id)
		ResponseError(c, errors.New("Convert workflow id error"), http.StatusBadRequest)
		return
	}

	workflowService := services.GetWorkflowService()
	workflow, err := workflowService.GetWorkflowByID(c, id)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, workflow)
	c.Abort()
	return
}

// HandlePUTWorkflowByID handle update workflow
// @Summary Update workflow by id
// @Description Update workflow by id
// @Produce json
// @Param workflow_id 	path 	int 				true 	"workflow id"
// @Param updateRequest body 	forms.WorkflowForm	true	"Update info"
// @Success 200 {object} forms.WorkflowDto "Workflow updated info"
// @Tags workflow
// @Router /workflows/:workflow_id [PUT]
func HandlePUTWorkflowByID(c *gin.Context) {
	workflow_id := c.Param("workflow_id")
	id, err := uuid.Parse(workflow_id)
	if err != nil {
		logger.Errorf("Convert workflow id failed: %s", workflow_id)
		ResponseError(c, errors.New("Convert workflow id error"), http.StatusBadRequest)
		return
	}

	var workflowForm forms.WorkflowForm
	err = c.Bind(&workflowForm)
	if err != nil {
		logger.Errorf("Error when parse workflow form: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	workflowService := services.GetWorkflowService()
	workflowResp, err := workflowService.UpdateWorkflow(c, id, workflowForm)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, workflowResp)
	return
}

// HandleDeleteWorkflow handle delete a workflow
// @Summary Delete workflow by id
// @Description Delete workflow by id
// @Produce json
// @Param workflow_id path int true "workflow id"
// @Success 200
// @Tags workflow
// @Router /workflows/:workflow_id [DELETE]
func HandleDeleteWorkflow(c *gin.Context) {
	workflow_id := c.Param("workflow_id")
	id, err := uuid.Parse(workflow_id)
	if err != nil {
		logger.Errorf("Convert workflow id failed: %s", workflow_id)
		ResponseError(c, errors.New("Convert workflow id error"), http.StatusBadRequest)
		return
	}

	workflowService := services.GetWorkflowService()
	err = workflowService.DeleteWorkflow(c, id)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// HandleGetRunsOfWorkflow return all runs belong a workflow
// @Summary Get run created from workflow id
// @Description Get run created from workflow id
// @Produce json
// @Param workflow_id path int true "workflow id"
// @Success 200 {object} forms.ListRunsDto "Run info"
// @Tags run
// @Router /workflows/:workflow_id/runs [GET]
func HandleGetRunsOfWorkflow(c *gin.Context) {
	pageSize, pageToken, filter, err := getFilterParam(c)
	if err != nil {
		return
	}

	workflow_id := c.Param("workflow_id")
	id, err := uuid.Parse(workflow_id)
	if err != nil {
		logger.Errorf("Convert workflow id failed: %s", workflow_id)
		ResponseError(c, errors.New("Convert workflow id error"), http.StatusBadRequest)
		return
	}

	//ctx := c
	//username := ctx.Value("UserName").(string)
	username := "tungnt99"
	workflowService := services.GetWorkflowService()
	runs, total, err := workflowService.GetWorkflowRuns(c, username, id, pageSize, pageToken, filter)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &forms.ListRunsDto{
		Runs:          runs,
		Total:         total,
		NextPageToken: strconv.Itoa(pageToken + 1),
	})
}
