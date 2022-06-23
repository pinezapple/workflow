package services

import (
	"context"

	"workflow/heimdall/repository"
	"workflow/heimdall/repository/entity"
	"workflow/heimdall/webserver/forms"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProjectService interface
type ProjectService interface {

	// GetProjects		   : return list projects based on filter
	GetProjects(ctx *gin.Context, pageSize int, pageToken int,
		filter map[string][]string) ([]forms.ProjectDto, int64, error)
	// GetProjectsFromAuth : returns list projects with matching auth resource path
	GetProjectsFromAuth(ctx context.Context, authPaths []string) (projects []forms.ProjectDto, err error)
	// GetProject	   : return Project based on uuid
	GetProject(ctx *gin.Context, id uuid.UUID) (forms.ProjectDto, error)
	// CreateProject	   : creates a project and return newly created project
	CreateProject(ctx *gin.Context, project forms.ProjectForm) (forms.ProjectDto, error)
	// UpdateProject	   : updates a project and return newly updated project
	UpdateProject(ctx *gin.Context, id uuid.UUID, project forms.ProjectForm) (forms.ProjectDto, error)
	// DeleteProject	   : deletes a project based on uuid
	DeleteProject(ctx *gin.Context, id uuid.UUID) error
	// GetProjectWorkflows : return list workflows of a project
	GetProjectWorkflows(ctx *gin.Context, username string, id uuid.UUID,
		pageSize int, pageToken int, filters map[string][]string) ([]forms.WorkflowDto, int64, error)

	CreateProjectFolder(ctx *gin.Context, id uuid.UUID, folder forms.FolderCreate) (forms.FolderDto, error)
	UpdateProjectFolder(ctx *gin.Context, folder forms.FolderUpdate) (forms.FolderDto, error)
	DeleteProjectFolder(ctx *gin.Context, id uuid.UUID) error
}

// GetProjectService return project service implement
func GetProjectService() ProjectService {
	return projectServiceImpl{}
}

type projectServiceImpl struct{}

func (service projectServiceImpl) GetProjects(ctx *gin.Context, pageSize int, pageToken int, filter map[string][]string) (dto []forms.ProjectDto, total int64, err error) {
	listProjects, total, err := repository.GetDAO().GetProjects(ctx, pageSize, pageToken, filter)
	if err != nil {
		return
	}
	for _, projectEntity := range listProjects {
		item := convertProjectEntity2Dto(&projectEntity)
		dto = append(dto, item)
	}
	return
}

func (service projectServiceImpl) GetProjectsFromAuth(ctx context.Context, authPaths []string) (projects []forms.ProjectDto, err error) {
	prjs, err := repository.GetDAO().GetAllProjectsFromAuth(ctx, authPaths)
	if err != nil {
		return nil, err
	}

	for i := range prjs {
		projects = append(projects, convertProjectEntity2Dto(&prjs[i]))
	}

	return projects, nil
}

func (service projectServiceImpl) GetProject(ctx *gin.Context, id uuid.UUID) (dto forms.ProjectDto, err error) {
	projectDB, err := repository.GetDAO().GetProject(ctx, id)
	if err != nil {
		return
	}
	dto = convertProjectEntity2Dto(&projectDB)
	return
}

func (service projectServiceImpl) CreateProject(ctx *gin.Context,
	project forms.ProjectForm) (forms.ProjectDto, error) {

	projectEntity, err := convertProjectForm2Entity(ctx, project, true)
	if err != nil {
		return forms.ProjectDto{}, err
	}

	dbDAO := repository.GetDAO()
	err = dbDAO.CreateProject(ctx, projectEntity)
	if err != nil {
		return forms.ProjectDto{}, err
	}

	return convertProjectEntity2Dto(projectEntity), nil
}

func (service projectServiceImpl) UpdateProject(ctx *gin.Context, id uuid.UUID, project forms.ProjectForm) (forms.ProjectDto, error) {
	updateProject, err := convertProjectForm2Entity(ctx, project, false)
	if err != nil {
		return forms.ProjectDto{}, err
	}

	updateProject.ID = id
	dbDAO := repository.GetDAO()
	err = dbDAO.UpdateProject(ctx, id, updateProject)
	if err != nil {
		return forms.ProjectDto{}, err
	}

	updatedProject, err := dbDAO.GetProject(ctx, id)
	if err != nil {
		logger.Errorf("Get project error: %s", id)
		return forms.ProjectDto{}, err
	}
	return convertProjectEntity2Dto(&updatedProject), nil
}

func (service projectServiceImpl) DeleteProject(ctx *gin.Context, id uuid.UUID) error {

	dbDAO := repository.GetDAO()
	if err := dbDAO.DeleteProject(ctx, id); err != nil {
		return err
	}

	return nil
}

func (service projectServiceImpl) GetProjectWorkflows(ctx *gin.Context, username string, id uuid.UUID, pageSize int, pageToken int,
	filters map[string][]string) (workflows []forms.WorkflowDto, total int64, err error) {
	workflowsEntity, total, err := repository.GetDAO().GetWorkflowsByProjectID(ctx, id, pageSize, pageToken, filters)
	if err != nil {
		logger.Errorf("Error when get project 's workflow: %s", err.Error())
		return
	}
	for _, workflowEntity := range workflowsEntity {
		var workflowDto = convertWorkflowEntity2Dto(&workflowEntity)
		workflows = append(workflows, workflowDto)
	}
	return
}

func (service projectServiceImpl) CreateProjectFolder(ctx *gin.Context, id uuid.UUID, folder forms.FolderCreate) (folderDto forms.FolderDto, err error) {
	var folderEntity = entity.FolderEntity{
		Name:      folder.Name,
		Path:      folder.Path,
		Author:    "tungnt99",
		ProjectID: id,
	}

	dbDAO := repository.GetDAO()
	err = dbDAO.AddProjectFolder(ctx, &folderEntity)

	if err != nil {
		return folderDto, err
	}

	folderDto = forms.FolderDto{
		ID:        folderEntity.ID,
		Name:      folderEntity.Name,
		Path:      folderEntity.Path,
		Author:    folderEntity.Author,
		CreatedAt: folderEntity.CreatedAt,
		UpdatedAt: folderEntity.UpdatedAt,
	}

	return folderDto, nil
}

func (service projectServiceImpl) UpdateProjectFolder(ctx *gin.Context, folder forms.FolderUpdate) (folderDto forms.FolderDto, err error) {
	var folderEntity = entity.FolderEntity{
		ID:   folder.ID,
		Name: folder.Name,
		Path: folder.Path,
	}

	dbDAO := repository.GetDAO()
	err = dbDAO.UpdateProjectFolder(ctx, &folderEntity)
	if err != nil {
		return folderDto, err
	}

	folderDto = forms.FolderDto{
		ID:        folderEntity.ID,
		Name:      folderEntity.Name,
		Path:      folderEntity.Path,
		Author:    folderEntity.Author,
		CreatedAt: folderEntity.CreatedAt,
		UpdatedAt: folderEntity.UpdatedAt,
	}

	return folderDto, nil
}

func (service projectServiceImpl) DeleteProjectFolder(ctx *gin.Context, id uuid.UUID) (err error) {
	dbDAO := repository.GetDAO()
	err = dbDAO.DeleteProjectFolder(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func convertProjectEntity2Dto(projectEntity *entity.ProjectEntity) (projectDto forms.ProjectDto) {
	projectDto = forms.ProjectDto{
		ID:          projectEntity.ID,
		Name:        projectEntity.Name,
		Description: projectEntity.Description,
		Summary:     projectEntity.Summary,
		CreatedAt:   projectEntity.CreatedAt,
		UpdatedAt:   projectEntity.UpdatedAt.Time,
		Author:      projectEntity.Author,
	}

	var folders []forms.FolderDto
	for _, folder := range projectEntity.Folders {
		folderDto := forms.FolderDto{
			ID:        folder.ID,
			Name:      folder.Name,
			Path:      folder.Path,
			Author:    folder.Author,
			CreatedAt: folder.CreatedAt,
			UpdatedAt: folder.UpdatedAt,
		}

		folders = append(folders, folderDto)
	}

	projectDto.Folders = folders

	return
}

func convertProjectForm2Entity(ctx *gin.Context, project forms.ProjectForm,
	generateID bool) (projectEntity *entity.ProjectEntity, err error) {
	projectEntity = &entity.ProjectEntity{
		Name:        project.Name,
		Description: project.Description,
		Summary:     project.Summary,
		Author:      "tungnt99",
	}
	if generateID {
		projectEntity.ID = uuid.New()
	}

	return
}
