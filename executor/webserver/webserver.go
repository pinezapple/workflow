package webserver

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"workflow/workflow-utils/model"
	"workflow/executor/core"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// WebServer booting web server by configuration
func WebServer(ctx context.Context) (fn model.Daemon, err error) {
	// get configs
	mainConfig := core.GetMainConfig()
	//lg := core.GetLogger()

	lg := core.GetLogger()
	// create admin by default
	e, server := echo.New(), &http.Server{
		Addr: fmt.Sprintf(":%d", mainConfig.HttpServerConfig.Port),
	}

	// Disable echo logging
	e.Logger.SetOutput(ioutil.Discard)

	// Recover middleware recovers from panics anywhere in the chain, prints stack trace and handles the control to the centralized HTTPErrorHandler.
	// Default stack size for trace is 4KB. For more example, please refer to https://echo.labstack.com/
	e.Use(mw.Recover())

	// Remove trailing slash middleware removes a trailing slash from the request URI.
	e.Pre(mw.RemoveTrailingSlash())

	// Set BodyLimit Middleware. It will panic if fail. For more example, please refer to https://echo.labstack.com/
	//	e.Use(mw.BodyLimit(mainConf.WebServer.BodyLimit))

	// Secure middleware provides protection against cross-site scripting (XSS) attack, content type sniffing, clickjacking, insecure connection and other code injection attacks.
	// For more example, please refer to https://echo.labstack.com/
	e.Use(mw.Secure())
	//	e.Use(mw.CORS())

	// TODO: init router
	//initRouter(e)

	lg.Info("HTTP Server is starting on " + strconv.Itoa(mainConfig.HttpServerConfig.Port))

	// Start server
	go func() {
		if err := e.StartServer(server); err != nil {
			lg.Errorf(err.Error())
		}
	}()

	fn = func() {
		<-ctx.Done()

		// try to shutdown server
		if err := e.Shutdown(context.Background()); err != nil {
			lg.Errorf(err.Error())
		} else {
			lg.Info("gracefully shutdown webserver")
		}
	}
	return
}
