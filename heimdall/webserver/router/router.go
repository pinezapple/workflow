// Package router package provide routing for the app and return a server to start
package router

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"workflow/heimdall/core"
	_ "workflow/heimdall/docs"
	"workflow/heimdall/webserver/controllers"
	"workflow/heimdall/webserver/middlewares"
)

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		str := fmt.Sprintf("IP: %v | Path: %v | Method: %v | Latency: %v | Status: %v", c.ClientIP(), c.Request.URL.Path+"?"+c.Request.URL.RawQuery, c.Request.Method, time.Since(t).Seconds(), c.Writer.Status())
		if !(c.Request.URL.Path == "/health" && c.Writer.Status() == 200) {
			core.GetLogger().Info(str)
		}
	}
}

//NewRouter definitions for services
func NewRouter() (router *gin.Engine) {
	router = gin.New()
	router.Use(gin.Recovery())
	router.Use(logRequest())

	//CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin,DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// The url pointing to API definition
	url := ginSwagger.URL("doc.json")
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Workflow handlers
	workflows := router.Group("/workflows", middlewares.DecodeToken())
	{
		workflows.GET("", controllers.HandleGETWorkflows)
		workflows.GET("/:workflow_id", controllers.HandleGETWorkflowByID)
		workflows.POST("", controllers.HandlePOSTWorkflow)
		workflows.PUT("/:workflow_id", controllers.HandlePUTWorkflowByID)
		workflows.DELETE("/:workflow_id", controllers.HandleDeleteWorkflow)
		workflows.GET("/:workflow_id/runs", controllers.HandleGetRunsOfWorkflow)
	}

	// Workflow runs handlers
	runs := router.Group("/runs", middlewares.DecodeToken())
	{
		runs.GET("", controllers.HandleGETRuns)
		runs.GET("/:run_id", controllers.HandleGETRun)
		runs.GET("/:run_id/status", controllers.HandleGETRunStatus)
		runs.POST("", controllers.HandlePOSTRun)
		// runs.PUT("/:run_id", controllers.HandlePUTRuns)
		// runs.POST("/:run_id/cancel", controllers.HandlePOSTCancelRunByID)
	}

	// Task handlers
	tasks := router.Group("/tasks", middlewares.DecodeToken())
	{
		tasks.GET("", controllers.HandleGETTasks)
		tasks.GET("/:task_id", controllers.HandleGETTask)
		tasks.POST("", controllers.HandlePOSTTask)
		// tasks.POST("/:id/cancel", controllers.HandlePOSTCancelTaskByID)
	}

	// Project handlers
	projects := router.Group("/projects", middlewares.DecodeToken())
	{
		projects.GET("", controllers.HandleGETProjects)
		projects.GET("/:project_id", controllers.HandleGETProject)
		projects.POST("", controllers.HandlePOSTProject)
		projects.PUT("/:project_id", controllers.HandlePUTProjectByID)
		projects.DELETE("/:project_id", controllers.HandleDeleteProjectByID)
		projects.GET("/:project_id/workflows", controllers.HandleGETWorkflowsOfProject)

		projects.POST("/:project_id/folders", controllers.HandlePOSTFolder)
		projects.PUT("/:project_id/folders", controllers.HandlePUTFolder)
		projects.DELETE("/:project_id/folders/:folder_id", controllers.HandleDELETEFolder)
	}

	// Task updater
	taskUpdate := router.Group("/internal/tasks")
	taskUpdaterCtrl := new(controllers.TaskUpdater)
	{
		taskUpdate.POST("/status", taskUpdaterCtrl.UpdateTaskStatus)
		taskUpdate.POST("/params", taskUpdaterCtrl.UpdateParamWithRegexAndFileSize)
		taskUpdate.POST("/output", taskUpdaterCtrl.UpdateTaskOutputLocation)
	}

	return router
}
