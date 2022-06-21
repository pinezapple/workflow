package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo"

	"workflow/executor/core"
	executorModel "workflow/executor/model"
)

func deleteK8STaskAsync(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *core.LogFormat, logResponse bool, err error) {
	req := request.(*executorModel.DeleteTaskK8SReq)

	//logger init
	lg = core.GetLogger()
	lg.Source = c.Request().RemoteAddr + "Delete jobs ASYNC"
	lg.Data = req.TaskID

	go deleteK8STask(req.TaskID)

	return http.StatusOK, data, lg, false, nil
}

func deleteK8STaskSync(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *core.LogFormat, logResponse bool, err error) {
	req := request.(*executorModel.DeleteTaskK8SReq)

	//logger init
	lg = core.GetLogger()
	lg.Source = c.Request().RemoteAddr + "Delete jobs SYNC"
	lg.Data = req.TaskID

	err = deleteK8STask(req.TaskID)
	if err != nil {
		return http.StatusInternalServerError, nil, lg, false, err
	}

	return http.StatusOK, data, lg, false, nil
}

func deleteK8STask(jobid []string) (err error) {
	for i := 0; i < len(jobid); i++ {
		err, _ = core.DeleteK8SJob(context.Background(), jobid[i], true)
		if err != nil {
			return err
		}
	}
	return nil
}
