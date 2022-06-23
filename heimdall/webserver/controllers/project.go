package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"workflow/heimdall/services"
	"workflow/heimdall/webserver/forms"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandlePOSTProject handle create project
// @Summary Create project
// @Description Create project
// @Produce json
// @Param sample body forms.ProjectForm true "Project info"
// @Success 200 {object} forms.ProjectDto "Create Project ok"
// @Tags project
// @Router /projects [POST]
func HandlePOSTProject(c *gin.Context) {
	var projectForm forms.ProjectForm
	err := c.Bind(&projectForm)
	if err != nil {
		logger.Errorf("Parse project error: %s", err.Error())
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	projectService := services.GetProjectService()
	project, err := projectService.CreateProject(c, projectForm)
	if err != nil {
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, project)
	return
}

// HandleGETProject handle get project by ID
// @Summary Get project
// @Description Get project
// @Produce json
// @Param project_id path string true "project ID"
// @Success 200 {object} forms.ProjectDto "Get Project ok"
// @Tags project
// @Router /projects/{project_id} [GET]
func HandleGETProject(c *gin.Context) {
	project_id := c.Param("project_id")
	id, err := uuid.Parse(project_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", project_id)
		return
	}

	projectService := services.GetProjectService()
	project, err := projectService.GetProject(c, id)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, project)
	c.Abort()
	return
}

// HandleGETProjects handle get projects
// @Summary 		Get projects
// @Description 	Get projects
// @Produce 		json
// @Param 			page_size 		query 	int 	false 	"Number of result"
// @Param 			page_token 		query 	int 	false 	"Current page"
// @Param 			filter 			query 	string 	false 	"Filter"
// @Param 			share			query	bool	false	"Get share projects"
// @Success 		200 			{array} forms.ProjectsDto 	"List Projects"
// @Tags 			project
// @Router 			/projects [GET]
func HandleGETProjects(c *gin.Context) {
	projectService := services.GetProjectService()
	pageSize, pageToken, filterMap, err := getFilterParam(c)
	if err != nil {
		return
	}

	projects, total, err := projectService.GetProjects(c, pageSize, pageToken, filterMap)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	response := &forms.ProjectsDto{
		Projects:      projects,
		NextPageToken: strconv.Itoa(pageToken + 1),
		Total:         total,
	}

	c.JSON(http.StatusOK, response)
	return
}

// HandleGETWorkflowsOfProject return all workflows belong a project
// @Summary Get Workflows of Project
// @Description Return all workflows belong a project
// @Produce json
// @Param project_uuid path string true "project UUID"
// @Param username query string false "username"
// @Param page_size query int false "Number of result"
// @Param page_token query int false "Current page"
// @Param filter query string false "Filter"
// @Success 200 {array} forms.WorkflowsDto "List workflows of project"
// @Tags project
// @Router /projects/{project_uuid}/workflows [GET]
func HandleGETWorkflowsOfProject(c *gin.Context) {
	pageSize, pageToken, filter, err := getFilterParam(c)
	if err != nil {
		return
	}

	project_id := c.Param("project_id")
	id, err := uuid.Parse(project_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", project_id)
		ResponseError(c, errors.New("Convert project id error"), http.StatusBadRequest)
		return
	}

	//ctx := c
	//username := ctx.Value("UserName").(string)
	username := "tungnt99"
	projectService := services.GetProjectService()
	workflows, total, err := projectService.GetProjectWorkflows(c, username, id, pageSize, pageToken, filter)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, &forms.WorkflowsDto{
		Workflows:     workflows,
		Total:         total,
		NextPageToken: strconv.Itoa(pageToken + 1),
	})
}

// HandlePUTProject handle update project
// @Summary Update project
// @Description Update project
// @Produce json
// @Param project_id path string true "project ID"
// @Param sample body forms.ProjectForm true "Project info"
// @Success 200 {object} forms.ProjectDto "Update Project ok"
// @Tags project
// @Router /projects/{project_id} [PUT]
func HandlePUTProjectByID(c *gin.Context) {
	project_id := c.Param("project_id")
	id, err := uuid.Parse(project_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", project_id)
		ResponseError(c, errors.New("Convert project id error"), http.StatusBadRequest)
		return
	}

	var projectForm forms.ProjectForm
	err = c.Bind(&projectForm)
	if err != nil {
		logger.Errorf("Error when parse project form: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	projectService := services.GetProjectService()
	project, err := projectService.UpdateProject(c, id, projectForm)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, project)
	return
}

// HandleDeleteProjectByID handle delete project
// @Summary Delete project
// @Description Delete project
// @Produce json
// @Param project_id path string true "project ID"
// @Success 200 {string} string "Delete Project ok"
// @Tags project
// @Router /projects/{project_id} [DELETE]
func HandleDeleteProjectByID(c *gin.Context) {
	project_id := c.Param("project_id")
	id, err := uuid.Parse(project_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", project_id)
		ResponseError(c, errors.New("Convert project id error"), http.StatusBadRequest)
		return
	}

	projectService := services.GetProjectService()
	err = projectService.DeleteProject(c, id)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// @Summary Create new folder in a project
// @Description Create new folder in a project
// @Produce json
// @Param project_id path string true "Project ID"
// @Param folder body forms.FolderCreate true "Folder info"
// @Success 200 {string} string "OK"
// @Tags project
// @Router /projects/{project_id}/folders [POST]
func HandlePOSTFolder(c *gin.Context) {
	project_id := c.Param("project_id")
	id, err := uuid.Parse(project_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", project_id)
		ResponseError(c, errors.New("Convert project id error"), http.StatusBadRequest)
		return
	}

	var folder forms.FolderCreate
	err = c.Bind(&folder)
	if err != nil {
		logger.Errorf("Error when parse project folder: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	projectService := services.GetProjectService()
	folderDto, err := projectService.CreateProjectFolder(c, id, folder)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, folderDto)
	return
}

// @Summary Update folder in a project
// @Description Update folder in a project
// @Produce json
// @Param project_id path string true "Project ID"
// @Param folder body forms.FolderUpdate true "Folder info"
// @Success 200 {string} string "OK"
// @Tags project
// @Router /projects/{project_id}/folders [PUT]
func HandlePUTFolder(c *gin.Context) {
	project_id := c.Param("project_id")
	_, err := uuid.Parse(project_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", project_id)
		ResponseError(c, errors.New("Convert project id error"), http.StatusBadRequest)
		return
	}

	var folder forms.FolderUpdate
	err = c.Bind(&folder)
	if err != nil {
		logger.Errorf("Error when parse project folder: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	projectService := services.GetProjectService()
	folderDto, err := projectService.UpdateProjectFolder(c, folder)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, folderDto)
	return
}

// @Summary Delete folder in a project
// @Description Delete folder in a project
// @Produce json
// @Param project_id path string true "Project ID"
// @Param folder_id path string true "Folder ID"
// @Success 200 {string} string "OK"
// @Tags project
// @Router /projects/{project_id}/folders/{folder_id} [DELETE]
func HandleDELETEFolder(c *gin.Context) {
	project_id := c.Param("project_id")
	_, err := uuid.Parse(project_id)
	if err != nil {
		logger.Errorf("Convert project id failed: %s", project_id)
		ResponseError(c, errors.New("Convert project id error"), http.StatusBadRequest)
		return
	}

	folder_id, err := uuid.Parse(c.Param("folder_id"))
	if err != nil {
		logger.Errorf("Convert folder id failed: %s", project_id)
		ResponseError(c, errors.New("Convert folder id error"), http.StatusBadRequest)
		return
	}

	projectService := services.GetProjectService()
	err = projectService.DeleteProjectFolder(c, folder_id)
	if err != nil {
		ResponseError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}
