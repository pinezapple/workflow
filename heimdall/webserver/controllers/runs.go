package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"workflow/heimdall/services"
	"workflow/heimdall/webserver/forms"
)

// HandleGETRuns handle get run list
// @Summary Get list run
// @Description Get list run
// @Produce json
// @Param 	page_size 		query 	int 	true "page size"
// @Param 	page_token 		query 	int 	true "page token"
// @Param 	filter query 	query	string 	true "Filter. Split by ;"
// @Success 200 {object} forms.ListRunsDto "a lot of runs"
// @Tags run
// @Router /runs [GET]
func HandleGETRuns(c *gin.Context) {
	ok := AuthzRequest(c, "/workflow/run", "read", "heimdall")
	if ok == false {
		ResponseError(c, errors.New("You do not have permission"), http.StatusForbidden)
		return
	}

	pageSize, pageToken, filterMap, _ := getFilterParam(c)

	runs, total, err := services.GetRunService().GetRuns(c, pageSize, pageToken, filterMap)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &forms.ListRunsDto{
		Runs:          runs,
		NextPageToken: strconv.Itoa(pageToken + 1),
		Total:         total,
	})
	return
}

// HandleGETRun handle get run by id
// @Summary Get run by id
// @Description Get run by id
// @Produce json
// @Param run_id path string true "run uuid"
// @Success 200 {object} forms.RunDto "Run info"
// @Tags run
// @Router /runs/:run_id [GET]
func HandleGETRun(c *gin.Context) {
	ok := AuthzRequest(c, "/workflow/run", "read", "heimdall")
	if ok == false {
		ResponseError(c, errors.New("You do not have permission"), http.StatusForbidden)
		return
	}

	run_id := c.Param("run_id")
	id, err := uuid.Parse(run_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", run_id)
		ResponseError(c, errors.New("Convert run id error"), http.StatusBadRequest)
		return
	}

	run, er := services.GetRunService().GetRun(c, id)
	if er != nil {
		ResponseError(c, er, 500)
		return
	}

	c.JSON(200, run)
	return
}

// HandlePOSTRun start run
// @Summary Create run and start
// @Description Create run and start
// @Produce json
// @Param createReq body forms.WorkflowRunForm true "Create Run Info"
// @Success 200 {string} run_id=string "Run uuid"
// @Tags run
// @Router /runs [POST]
func HandlePOSTRun(c *gin.Context) {
	ok := AuthzRequest(c, "/workflow/run", "create", "heimdall")
	if ok == false {
		ResponseError(c, errors.New("You do not have permission"), http.StatusForbidden)
		return
	}

	var runForm *forms.WorkflowRunForm
	err := c.BindJSON(&runForm)
	if err != nil {
		logger.Errorf("Error when parse rundto: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	runService := services.GetRunService()
	runDto, err := runService.CreateRun(c, runForm)
	if err != nil {
		logger.Errorf("Error when create run: %s", err.Error())
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}
	logger.Infof("Run uuid: %s", runDto.ID)

	c.JSON(http.StatusOK, gin.H{
		"run_id": runDto.ID,
	})
}

// HandleGETRunStatus return run status
// @Summary Get run status by id
// @Description Get run status by id
// @Produce json
// @Param run_id path string true "run uuid"
// @Success 200 {object} forms.RunDto "Run info"
// @Tags run
// @Router /runs/:run_id/status [GET]
func HandleGETRunStatus(c *gin.Context) {
	ok := AuthzRequest(c, "/workflow/run", "read", "heimdall")
	if ok == false {
		ResponseError(c, errors.New("You do not have permission"), http.StatusForbidden)
		return
	}

	run_id := c.Param("run_id")
	id, err := uuid.Parse(run_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", run_id)
		ResponseError(c, errors.New("Convert run id error"), http.StatusBadRequest)
		return
	}

	run, er := services.GetRunService().GetRunStatus(c, id)
	if er != nil {
		ResponseError(c, er, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, run)
	return
}
