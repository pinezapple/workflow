package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"workflow/workflow-utils/model"
	"workflow/executor/core"

	"github.com/labstack/echo"
)

// ExecHandler execute handler
func ExecHandler(c echo.Context, expect interface{}, invoke func(c echo.Context, req interface{}) (int, interface{}, *core.LogFormat, bool, error)) error {
	var payload []byte
	var err error

	body := c.Request().Body
	defer func() {
		_ = body.Close()
	}()

	if payload, err = ioutil.ReadAll(body); err != nil {
		return c.JSON(http.StatusOK, &model.Response{
			Error: model.ResponseError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}})
	}

	if expect != nil { // parse req
		if err = json.Unmarshal(payload, expect); err != nil {
			return c.JSON(http.StatusOK, &model.Response{
				Error: model.ResponseError{
					Code:    http.StatusBadRequest,
					Message: core.ErrBadRequest.Error(),
				}})
		}
	}

	statusCode, data, lg, logRespData, err := invoke(c, expect)
	if err != nil {
		if lg != nil {
			if logRespData {
				lg.Success = data
			}
			lg.Errorf(err.Error())
		}
		return c.JSON(http.StatusOK, &model.Response{
			Error: model.ResponseError{
				Code:    statusCode,
				Message: err.Error(),
			},
		})
	}

	if lg != nil {
		if logRespData {
			lg.Success = data
		}
		lg.Info("Success")
	}

	return c.JSON(http.StatusOK, &model.Response{
		Data: data,
		Error: model.ResponseError{
			Code: http.StatusOK,
		},
	})
}
