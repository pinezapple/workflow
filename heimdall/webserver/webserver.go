// Package webserver  contains function to create a web server
package webserver

import (
	"context"
	"fmt"
	"net/http"

	"workflow/heimdall/core"
	_ "workflow/heimdall/docs"
	"workflow/heimdall/repository"
	"workflow/heimdall/webserver/router"
	"workflow/workflow-utils/model"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	logger *core.Logger = core.GetLogger()
)

func constructPostgresDSN(conf *core.DbConfig) (dsn string) {
	dsn = "user=" + conf.User + " password=" + conf.Password + " dbname=" + conf.DBName + " host=" + conf.Host + " port=" + conf.Port + " sslmode=disable"
	return
}

func initialDatabaseConnection() (err error) {
	mainConfig := core.GetConfig()
	dsn := constructPostgresDSN(mainConfig.DB)
	err = repository.GetDAO().InitDBConnection(dsn)
	return
}

func initialTableDatabase(ctx context.Context) error {
	return repository.GetDAO().AutoMigrate(ctx)
}

// InitDabase init database connection and migrate table
func InitDabase(ctx context.Context) {
	if err := initialDatabaseConnection(); err != nil {
		logger.Errorf("Init database connection error %s", err.Error())
		panic(err)
	} else {
		logger.Info("Init database connection success")
	}
	if err := initialTableDatabase(ctx); err != nil {
		logger.Errorf("Auto migrate table database error %s", err.Error())
		panic(err)
	} else {
		logger.Info("Auto migrate table database success")
	}
}

func WebServer(ctx context.Context) (fn model.Daemon, err error) {
	e := echo.New()

	e.GET("/docs/*", echoSwagger.WrapHandler)

	config := core.GetConfig()
	logger = core.InitLogger("heimdall", config.Server.LogLevel, config.Server.Environment)

	logger.Info("Starting webserver")

	InitDabase(ctx)
	r := router.NewRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.
			ErrServerClosed {
			logger.Fatalf("list: %s\n", err)
		}
	}()

	fn = func() {
		<-ctx.Done()

		// try to shutdown server
		if err := server.Shutdown(ctx); err != nil {
			logger.Errorf("Server shutdown: %v", err)
		} else {
			logger.Info("Gracefully shutdown webserver")
		}
	}

	return
}
