package handler

import (
	"github.com/labstack/echo"
	"github.com/vfluxus/workflow/executor/controller"

	executorModel "github.com/vfluxus/workflow/executor/model"
)

func DeleteK8STaskAsync(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &executorModel.DeleteTaskK8SReq{}, deleteK8STaskAsync)
}

func DeleteK8STaskSync(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &executorModel.DeleteTaskK8SReq{}, deleteK8STaskSync)
}
