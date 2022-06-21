package controllers

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"workflow/heimdall/core"
	"workflow/heimdall/repository/gormdb"
)

func TestMain(m *testing.M) {
	core.ReadConfig("/home/ubuntu/Desktop/cloneworkflow/heimdall/app/dev.yaml")
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable",
		core.GetConfig().DB.User, core.GetConfig().DB.Password, core.GetConfig().DB.DBName,
		core.GetConfig().DB.Port, core.GetConfig().DB.Host)
	gormdb.GetGormDB().InitDBConnection(dsn)
	router := gin.Default()
	router.GET("/workflows", HandleGETWorkflows)
	router.GET("/workflows/:workflow_id", HandleGETWorkflowByID)
	router.POST("/workflows", HandlePOSTWorkflow)
	router.PUT("/workflows/:workflow_id", HandlePUTWorkflowByID)
	router.DELETE("/workflows/:workflow_id", HandleDeleteWorkflow)
	router.GET("/workflows/:workflow_id/runs", HandleGetRunsOfWorkflow)
	router.GET("/runs", HandleGETRuns)
	router.POST("/runs", HandlePOSTRun)
	router.GET("/runs/:run_id", HandleGETRun)
	router.GET("/runs/:run_id/status", HandleGETRunStatus)
	//router.GET("/run/:run_id/status", HandleGETRunStatus)
	//router.POST("/runs", HandlePOSTRun)
	router.Run(":8084")
}
