package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/vfluxus/heimdall/repository"
	"github.com/vfluxus/heimdall/services"
	"github.com/vfluxus/heimdall/webserver/forms"

	"github.com/gin-gonic/gin"
)

type TaskUpdater struct {
}

// NOTE: FIXED TO 1 RESPONSE
func responseUpdateStatusError(c *gin.Context, code int, err error) {
	logger.Errorf("Can not update task status: %v", err)
	errorResponse := forms.GetTaskUpdateResp(http.StatusInternalServerError, err.Error(), nil)
	c.JSON(http.StatusInternalServerError, errorResponse)
	c.Abort()
}

// UpdateTaskStatus ....
// @Summary Update task status by scheduler
// @Description Update task status by scheduler
// @Produce json
// @Param updateRequest body forms.UpdateTaskStatusForm true "Update info"
// @Success 200 {object} forms.TaskUpdateResponse{data=string} "Update success"
// @Tags taskupdate
// @Router /internal/tasks/status [POST]
func (tU *TaskUpdater) UpdateTaskStatus(c *gin.Context) {
	var (
		updateForm = new(forms.UpdateTaskStatusForm)
	)

	if err := c.ShouldBindJSON(updateForm); err != nil {
		logger.Errorf("Can not bind json request: %v", err)
		ResponseError(c, err, http.StatusNotAcceptable)
		return
	}

	if data, err := json.MarshalIndent(updateForm, "", "    "); err == nil {
		logger.Debugf("\nUpdate Task status form: %v", string(data))
	}

	if err := services.GetTaskStatusUpdater().UpdateStatusAuto(c, updateForm.TaskIDs, updateForm.Status); err != nil {
		logger.Errorf("Update status error: %v", err)
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	logger.Infof("Update Task status OK: %v | Status: %v |", updateForm.TaskIDs, updateForm.Status)
	c.JSON(http.StatusOK, forms.GetTaskUpdateResp(http.StatusOK, "", nil))
}

// UpdateTaskOutputLocation ....
// @Summary Update task output location by scheduler
// @Description Update task output location by scheduler
// @Produce json
// @Param updateRequest body forms.UpdateTaskOutputLocationForm true "Update info"
// @Success 200 {object} forms.TaskUpdateResponse{data=string} "Update success"
// @Tags taskupdate
// @Router /internal/tasks/output [POST]
func (tU *TaskUpdater) UpdateTaskOutputLocation(c *gin.Context) {
	var (
		updateForm = new(forms.UpdateTaskOutputLocationForm)
	)

	if err := c.ShouldBindJSON(updateForm); err != nil {
		logger.Errorf("Can not bind json: %v", err)
		responseUpdateStatusError(c, http.StatusNotAcceptable, err)
		return
	}

	if data, err := json.MarshalIndent(updateForm, "", "    "); err == nil {
		logger.Debugf("\nUpdate TaskOutputLocation form: %v", string(data))
	}

	if err := repository.GetTaskDAO().UpdateOutputLocation(c, updateForm.TaskID, updateForm.OutputLocation); err != nil {
		logger.Errorf("Can not update output location: %v", err)
		responseUpdateStatusError(c, http.StatusInternalServerError, err)
		return
	}

	logger.Infof("Update output location OK. ID: %v", updateForm.TaskID)
	c.JSON(http.StatusOK, forms.GetTaskUpdateResp(http.StatusOK, "", nil))
}

// UpdateParamWithRegexAndFileSize ...
// @Summary Update ParamRegex And FileSize by scheduler
// @Description Update ParamRegex And FileSize by scheduler
// @Produce json
// @Param UpdateReq body forms.UpdateParamWithRegexFileSizeForm true "Update info"
// @Success 200 {object} forms.TaskUpdateResponse{data=string} "Update success"
// @Tags taskupdate
// @Router /internal/tasks/params [POST]
func (tU *TaskUpdater) UpdateParamWithRegexAndFileSize(c *gin.Context) {
	// NOTE: Change to 1 transaction in DB
	var (
		updateForm = new(forms.UpdateParamWithRegexFileSizeForm)
	)

	if err := c.ShouldBindJSON(updateForm); err != nil {
		logger.Errorf("Can not bind json: %v", err)
		responseUpdateStatusError(c, http.StatusNotAcceptable, err)
		return
	}

	if data, err := json.MarshalIndent(updateForm, "", "    "); err == nil {
		logger.Debugf("\nUpdate ParamWithRegex And FileSize form: %v", string(data))
	}

	if len(updateForm.FileSize) != 0 {
		if err := repository.GetTaskDAO().UpdateFileSize(c, updateForm.TaskID, updateForm.FileSize); err != nil {
			logger.Errorf("Can not update file size: %v", err)
			responseUpdateStatusError(c, http.StatusInternalServerError, err)
			return
		}
		logger.Infof("Update file size OK. ID: %v", updateForm.TaskID)
	}

	if len(updateForm.ParamWithRegex) != 0 {
		data, err := json.Marshal(updateForm.ParamWithRegex)
		if err != nil {
			logger.Errorf("Can not marshal param with regex: %v", err)
			responseUpdateStatusError(c, http.StatusInternalServerError, err)
			return
		}

		if err := repository.GetTaskDAO().UpdateParamWithRegex(c, updateForm.TaskID, data); err != nil {
			logger.Errorf("Can not update Param with regex: %v", err)
			responseUpdateStatusError(c, http.StatusInternalServerError, err)
			return
		}
		logger.Infof("Update param with regex ok. ID: %v", updateForm.TaskID)
	}

	c.JSON(http.StatusOK, forms.GetTaskUpdateResp(http.StatusOK, "", nil))
}
