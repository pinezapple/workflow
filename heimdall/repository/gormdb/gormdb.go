// Package gormdb manage database
package gormdb

import (
	"context"
	"fmt"

	"workflow/heimdall/core"
	"workflow/heimdall/repository/entity"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/clause"

	//"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	//"gorm.io/gorm/clause"
)

type gormHeimdall struct{}

var (
	gDB                 = &gorm.DB{}
	gor    gormHeimdall = gormHeimdall{}
	logger *core.Logger = core.GetLogger()
)

// GetGormDB returns GormDB
func GetGormDB() *gormHeimdall {
	return &gor
}

func (c *gormHeimdall) AutoMigrate(ctx context.Context) (err error) {
	err = gDB.WithContext(ctx).
		AutoMigrate(
			&entity.ProjectEntity{},
			&entity.WorkflowEntity{},
			&entity.WorkflowStepEntity{},
			&entity.RunEntity{},
			&entity.TaskEntity{},
			&entity.FolderEntity{})
	if err != nil {
		return err
	}
	return
}

func (c *gormHeimdall) InitDBConnection(dsn string) (err error) {
	gDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlog.Default.LogMode(gormlog.Error),
	})
	if err != nil {
		return
	}

	return nil
}

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- WORKFLOWS ----------------------------------------------------------

func (c *gormHeimdall) GetWorkflow(ctx context.Context, id uuid.UUID) (wf entity.WorkflowEntity, er error) {
	er = gDB.WithContext(ctx).Model(&entity.WorkflowEntity{}).Where(&entity.WorkflowEntity{ID: id}).Preload(clause.Associations).Find(&wf).Error
	if er != nil {
		return
	}
	return
}

func (c *gormHeimdall) GetWorkflows(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) (wfs []entity.WorkflowEntity, total int64, er error) {
	thisDB := gDB.WithContext(ctx).Model(&entity.WorkflowEntity{})
	for k, v := range filter {
		thisDB = thisDB.Where(fmt.Sprintf("%s %s ?", k, v[0]), v[1])
	}
	thisDB.Preload(clause.Associations).Count(&total)

	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Offset(offset).Limit(pageSize)
	er = thisDB.Preload(clause.Associations).Find(&wfs).Error
	if er != nil {
		return
	}
	return
}

func (c *gormHeimdall) CreateWorkflow(ctx context.Context, workflow *entity.WorkflowEntity) error {
	if err := gDB.WithContext(ctx).Model(&entity.WorkflowEntity{}).Create(workflow).Error; err != nil {
		logger.Errorf("Create workflow error: %s", err.Error())
		return err
	}

	return nil
}

func (c *gormHeimdall) UpdateWorkflow(ctx context.Context, id uuid.UUID, workflow *entity.WorkflowEntity) error {
	db := gDB.Begin()
	if err := db.WithContext(ctx).Where(&entity.WorkflowStepEntity{WorkflowID: id}).Delete(&entity.WorkflowStepEntity{}).Error; err != nil {
		logger.Errorf("Remove workflow steps error: %s", err.Error())
		return err
	}

	if err := db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Where(&entity.WorkflowEntity{ID: id}).Updates(workflow).Error; err != nil {
		logger.Errorf("Update workflow error: %s", err.Error())
		db.Rollback()
		return err
	}

	// if err := db.WithContext(ctx).CreateInBatches(workflow.Steps, len(workflow.Steps)).Error; err != nil {
	// 	logger.Errorf("Insert workflow steps error: %s", err.Error())
	// 	db.Rollback()
	// 	return err
	// }
	db.Commit()

	return nil
}

func (c *gormHeimdall) DeleteWorkflow(ctx context.Context, id uuid.UUID) error {
	gDB.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Where(&entity.WorkflowEntity{ID: id}).Delete(&entity.WorkflowEntity{})
	return nil
}

// -----------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- RUN ----------------------------------------------------------

func (c *gormHeimdall) GetRunsByWorkflowID(ctx context.Context, username string, id uuid.UUID, pageSize int, pageToken int, filter map[string][]string) (runs []entity.RunEntity, total int64, er error) {
	var workflow = entity.WorkflowEntity{}
	er = gDB.WithContext(ctx).Model(&entity.WorkflowEntity{}).Where(&entity.WorkflowEntity{ID: id}).First(&workflow).Error
	if er != nil {
		return
	}
	thisDB := gDB.WithContext(ctx).Model(&entity.RunEntity{}).Where(&entity.RunEntity{WorkflowID: workflow.ID, UserName: username})
	for k, v := range filter {
		thisDB = thisDB.Where(fmt.Sprintf("%s %s ?", k, v[0]), v[1])
	}
	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Count(&total).Offset(offset).Limit(pageSize).Order("created_at DESC")
	er = thisDB.Omit(clause.Associations).Find(&runs).Error
	if er != nil {
		return
	}
	return
}

func (c *gormHeimdall) GetRun(ctx context.Context, id uuid.UUID) (run entity.RunEntity, er error) {
	er = gDB.WithContext(ctx).Model(&entity.RunEntity{}).
		Where(&entity.RunEntity{ID: id}).
		Preload("Tasks", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("task_entities.started_time ASC")
		}).
		Find(&run).Error
	if er != nil {
		return
	}
	return
}

func (c *gormHeimdall) GetRuns(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) (runs []entity.RunEntity, total int64, er error) {
	thisDB := gDB.WithContext(ctx).Model(&entity.RunEntity{}).Order("created_at DESC, id")
	for k, v := range filter {
		thisDB = thisDB.Where(fmt.Sprintf("%s %s ?", k, v[0]), v[1])
	}
	thisDB.Omit(clause.Associations).Count(&total)

	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Offset(offset).Limit(pageSize)
	er = thisDB.Preload(clause.Associations).Find(&runs).Error
	if er != nil {
		return
	}
	return
}

func (c *gormHeimdall) CreateRun(ctx context.Context, runEntity *entity.RunEntity) error {
	if err := gDB.WithContext(ctx).Model(&entity.RunEntity{}).Create(&runEntity).Error; err != nil {
		logger.Errorf("Create run error: %v", err.Error())
		return err
	}

	return nil
}

func (c *gormHeimdall) UpdateRun(ctx context.Context, runEntity *entity.RunEntity) error {
	if err := gDB.WithContext(ctx).Model(runEntity).
		Where("id = ?", runEntity.ID).
		Updates(runEntity).Error; err != nil {
		logger.Errorf("Update run error: %v", err.Error())
		return err
	}
	return nil
}

func (c *gormHeimdall) CancelRun(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (c *gormHeimdall) DeleteRun(ctx context.Context, runEntity *entity.RunEntity) error {
	return gDB.WithContext(ctx).Delete(runEntity).Error
}

// -------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- TASKS ----------------------------------------------------------

func (c *gormHeimdall) GetTasks(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) (listTasks []entity.TaskEntity, total int64, err error) {
	// Special case when user wants to get all tasks belong to a run by send
	// request such as http://heimdall/tasks?run_uuid=d9921d7e-947c-4b0f-a527-58ceb302c836
	// Because task entity does not have run uuid so need to catch this query
	// separately
	var db *gorm.DB
	if runUUID, ok := filter["run_uuid"]; ok {
		logger.Infof("Current uuid %s", runUUID)
		db = gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Joins("JOIN run_entities on run_entities.id=? and run_entities.id = task_entities.run_id", runUUID[1])
		delete(filter, "run_uuid")
	} else {
		logger.Infof("Current filter %s", filter)
		db = gDB.WithContext(ctx).Model(&entity.TaskEntity{})
	}

	// FIXME: Where conditions are not included in query
	for k, v := range filter {
		db.Where(fmt.Sprintf("%s %s ?", k, v[0]), v[1])
	}
	db.Where("task_id NOT LIKE ? AND task_id NOT LIKE ?", "%bigbang%", "%ragnarok%")
	db.Order("started_time DESC")

	offset := (pageToken - 1) * pageSize
	db.Count(&total).Offset(offset).Limit(pageSize)
	err = db.Find(&listTasks).Error
	return
}

func (c *gormHeimdall) GetTask(ctx context.Context, id uuid.UUID) (task entity.TaskEntity, err error) {
	err = gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Where(&entity.TaskEntity{ID: id}).Preload(clause.Associations).Find(&task).Error
	return
}

func (c *gormHeimdall) CreateTask(ctx context.Context, taskEntity *entity.TaskEntity) error {
	if err := gDB.WithContext(ctx).Model(&entity.TaskEntity{}).Create(taskEntity).Error; err != nil {
		logger.Errorf("Create task error: %v", err.Error())
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- PROJECT ----------------------------------------------------------

func (c *gormHeimdall) GetProject(ctx context.Context, id uuid.UUID) (project entity.ProjectEntity, err error) {
	err = gDB.WithContext(ctx).Model(&entity.ProjectEntity{}).Where(&entity.ProjectEntity{ID: id}).Preload(clause.Associations).Find(&project).Error
	if err != nil {
		return
	}
	return
}

func (c *gormHeimdall) GetProjects(ctx context.Context, pageSize int, pageToken int, filter map[string][]string) (projects []entity.ProjectEntity, total int64, err error) {
	thisDB := gDB.WithContext(ctx).Model(&entity.ProjectEntity{})
	for k, v := range filter {
		thisDB = thisDB.Where(fmt.Sprintf("%s %s ?", k, v[0]), v[1])
	}
	thisDB.Omit(clause.Associations).Count(&total)

	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Offset(offset).Limit(pageSize)
	err = thisDB.Preload(clause.Associations).Find(&projects).Error
	if err != nil {
		return
	}
	return
}

func (c *gormHeimdall) GetProjectsFromAuth(ctx context.Context, pageSize int, pageToken int, authPath []string) ([]entity.ProjectEntity, int64, error) {
	var (
		prjs   []entity.ProjectEntity
		total  int64
		offset = (pageToken - 1) * pageSize
	)

	err := gDB.WithContext(ctx).Model(&entity.ProjectEntity{}).
		Omit(clause.Associations).
		Order("id ASC").
		Offset(offset).Limit(pageSize).
		Count(&total).Find(&prjs).Error

	if err != nil {
		return nil, 0, err
	}

	return prjs, total, nil
}

func (c *gormHeimdall) GetAllProjectsFromAuth(ctx context.Context, authPaths []string) (projects []entity.ProjectEntity, err error) {
	err = gDB.WithContext(ctx).Model(&entity.ProjectEntity{}).
		Preload(clause.Associations).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *gormHeimdall) CreateProject(ctx context.Context, projectEntity *entity.ProjectEntity) error {
	if err := gDB.WithContext(ctx).Model(&entity.ProjectEntity{}).Create(&projectEntity).Error; err != nil {
		logger.Errorf("Create project error: %s", err.Error())
		return err
	}

	return nil
}

func (c *gormHeimdall) UpdateProject(ctx context.Context, id uuid.UUID, project *entity.ProjectEntity) error {
	db := gDB.Begin()
	if err := db.WithContext(ctx).Omit(clause.Associations).Model(entity.ProjectEntity{}).Where(&entity.ProjectEntity{ID: id}).Updates(project).Error; err != nil {
		logger.Errorf("Update project error: %s", err.Error())
		db.Rollback()
		return err
	}
	db.Commit()

	return nil
}

func (c *gormHeimdall) DeleteProject(ctx context.Context, id uuid.UUID) error {
	gDB.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Where(&entity.ProjectEntity{ID: id}).Delete(&entity.ProjectEntity{})
	return nil
}

func (c *gormHeimdall) GetWorkflowsByProjectID(ctx context.Context, id uuid.UUID,
	pageSize int, pageToken int, filter map[string][]string) (workflows []entity.WorkflowEntity, total int64, err error) {

	var project = entity.ProjectEntity{}
	err = gDB.WithContext(ctx).Model(&entity.ProjectEntity{}).Where(&entity.ProjectEntity{ID: id}).First(&project).Error
	if err != nil {
		return
	}

	thisDB := gDB.WithContext(ctx).Model(&entity.WorkflowEntity{}).Where(&entity.WorkflowEntity{ProjectID: project.ID})
	for k, v := range filter {
		thisDB = thisDB.Where(fmt.Sprintf("%s %s ?", k, v[0]), v[1])
	}
	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Count(&total).Offset(offset).Limit(pageSize).Order("created_at DESC")
	err = thisDB.Preload("Project").Find(&workflows).Error
	if err != nil {
		return
	}
	return
}

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FOLDER ----------------------------------------------------------

func (c *gormHeimdall) AddProjectFolder(ctx context.Context, folder *entity.FolderEntity) (err error) {
	err = gDB.WithContext(ctx).Create(folder).Error
	if err != nil {
		logger.Errorf("Add project folder error: %s", err.Error())
	}
	return
}

func (c *gormHeimdall) UpdateProjectFolder(ctx context.Context, folder *entity.FolderEntity) (err error) {
	err = gDB.WithContext(ctx).Updates(folder).Error
	if err != nil {
		logger.Errorf("Update project folder error: %s", err.Error())
	}
	return
}

func (c *gormHeimdall) DeleteProjectFolder(ctx context.Context, folder_id uuid.UUID) (err error) {
	err = gDB.WithContext(ctx).Delete(&entity.FolderEntity{}, folder_id).Error
	if err != nil {
		logger.Errorf("Delete project folder error: %s", err.Error())
	}
	return
}
