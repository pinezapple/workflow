package webserver

import (
	"context"
	"fmt"
	"net/http"

	"strconv"
	"time"

	"workflow/valkyrie/controller/handler"
	"workflow/valkyrie/core"
	"workflow/valkyrie/middleware"
	"workflow/workflow-utils/model"

	echo "github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"

	// echopprof "github.com/sevenNt/echo-pprof"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "workflow/valkyrie/docs"
)

func logRequest(logger *core.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := time.Now()

			str := fmt.Sprintf("%v %v", c.Request().Method, c.Request().URL.Path+"?"+c.Request().URL.RawQuery)
			logger.Info(str)

			err := next(c)

			str = fmt.Sprintf("%v %v %v %v", c.Request().Method, c.Response().Status, time.Since(t).Seconds(), c.Request().URL.Path+"?"+c.Request().URL.RawQuery)
			logger.Info(str)

			return err
		}
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func WebServer(ctx context.Context) (fn model.Daemon, err error) {
	// get configs
	mainConfig := core.GetMainConfig()
	logger := core.InitLogger(mainConfig.ServiceName, mainConfig.LogLevel, mainConfig.Environment)

	// create admin by default
	e, server := echo.New(), &http.Server{
		Addr:         fmt.Sprintf(":%d", mainConfig.HttpServerConfig.Port),
		IdleTimeout:  2 * time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 60 * time.Minute,
		// ReadHeaderTimeout: 1 * time.Minute,
	}

	// Disable echo logging
	// e.Logger.SetOutput(ioutil.Discard)

	// Recover middleware recovers from panics anywhere in the chain, prints stack trace and handles the control to the centralized HTTPErrorHandler.
	// Default stack size for trace is 4KB. For more example, please refer to https://echo.labstack.com/
	//e.Use(mw.Recover())
	if mainConfig.Environment == "DEVELOPMENT" {
		e.Use(mw.CORSWithConfig(mw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowMethods: []string{echo.POST, echo.PUT, echo.OPTIONS, echo.PATCH},
		}))
	}

	// Remove trailing slash middleware removes a trailing slash from the request URI.
	e.Pre(mw.RemoveTrailingSlash())

	// Set BodyLimit Middleware. It will panic if fail. For more example, please refer to https://echo.labstack.com/
	//	e.Use(mw.BodyLimit(mainConf.WebServer.BodyLimit))

	// Secure middleware provides protection against cross-site scripting (XSS) attack, content type sniffing, clickjacking, insecure connection and other code injection attacks.
	// For more example, please refer to https://echo.labstack.com/
	e.Use(mw.Secure())
	e.Use(logRequest(logger))

	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	// echopprof.Wrap(e)

	if mainConfig.HardDiskOnly {
		initHardDiskOnlyRouter(e)
	} else {
		initRouter(e)
	}

	logger.Info("HTTP Server is starting on " + strconv.Itoa(mainConfig.HttpServerConfig.Port))

	// Start server
	go func() {
		if err := e.StartServer(server); err != nil {
			logger.Error(err.Error())
		}
	}()

	fn = func() {
		<-ctx.Done()

		// try to shutdown server
		if err := e.Shutdown(context.Background()); err != nil {
			logger.Error(err.Error())
		} else {
			logger.Error("gracefully shutdown webserver")
		}
	}
	return
}

func initHardDiskOnlyRouter(e *echo.Echo) {
	files := e.Group("/files", middleware.DecodeJWT)
	{
		files.GET("", handler.GetFiles)
		files.GET("/uploaded", handler.GetUserUploadFiles)
		files.GET("/generated", handler.GetGeneratedFiles)
		//files.POST("/upload", handler.UserUploadFileToMinio)
		//files.GET("/sample/:sample_name", handler.CheckExistSample)

		// TODO(tuandn8) Maybe error when upload time longer than expired time of access token
		files.GET("/_resumable", handler.GetResumableUploadFile)
		files.POST("/_resumable", handler.ResumableUploadFile)

		files.GET("/:file_name", handler.DownloadFileFromHardDisk)
		files.PUT("/project_path", handler.UpdateFileProjectPath)

		files.DELETE("/file/:file_id", handler.DeleteFileFromHardDisk)
	}

	internal := e.Group("/internal")
	{
		internal.POST("/files", handler.HandleGeneratedFile)
		internal.DELETE("/files", handler.DeleteFailRunDir)
	}

	dataset := e.Group("/datasets", middleware.DecodeJWT)
	{
		dataset.GET("", handler.GetDatasetFilter)
		dataset.POST("", handler.CreateDataset)
	}

	sample := e.Group("/samples", middleware.DecodeJWT)
	{
		sample.GET("", handler.GetSampleFilter)
		sample.POST("", handler.CreateSample)
		sample.GET("/name/:sample_name", handler.CheckExistSample)
		sample.GET("/workflow/:workflow_uuid", handler.GetSamplesByWorkflow)
		sample.GET("/dataset/:dataset_uuid", handler.GetSampleByDataset)
		sample.GET("/:sample_uuid", handler.GetSampleDetail)
	}

}

func initRouter(e *echo.Echo) {
	files := e.Group("/files", middleware.DecodeJWT)
	{
		//files.GET("/", handler.GetUserFilesWithFilters)
		//files.POST("/upload", handler.UserUploadFileToMinio)
		//files.GET("/sample/:sample_name", handler.CheckExistSample)
		files.POST("/reupload", handler.ReuploadFileToMinio)
		files.POST("/presigned", handler.GetPresignedDownloadURL)
	}

	internal := e.Group("/internal")
	{
		internal.POST("/files/upload", handler.HandleGeneratedFile)
		internal.POST("/files/delete", handler.DeleteFailRunDir)
	}
}
