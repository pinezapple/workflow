package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/vfluxus/valkyrie/core"
	"github.com/vfluxus/workflow-utils/model"
)

var (
	logger *core.Logger = core.GetLogger()
)

// ExecHandler execute handler
func ExecHandler(c echo.Context, expect interface{}, invoke func(c echo.Context, req interface{}) (int, interface{}, error)) error {
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

	statusCode, data, err := invoke(c, expect)
	if err != nil {
		return c.JSON(http.StatusOK, &model.Response{
			Error: model.ResponseError{
				Code:    statusCode,
				Message: err.Error(),
			},
		})
	}

	logger.Info("Success")

	return c.JSON(http.StatusOK, &model.Response{
		Data: data,
		Error: model.ResponseError{
			Code: http.StatusOK,
		},
	})
}

// ExecHandler execute handler
func ExecHandlerForUploads(c echo.Context, expect interface{}, invoke func(c echo.Context, req interface{}) (int, interface{}, error)) error {
	statusCode, data, err := invoke(c, expect)
	if err != nil {
		return c.JSON(http.StatusOK, &model.Response{
			Error: model.ResponseError{
				Code:    statusCode,
				Message: err.Error(),
			},
		})
	}

	logger.Info("Success")

	return c.JSON(http.StatusOK, &model.Response{
		Data: data,
		Error: model.ResponseError{
			Code: http.StatusOK,
		},
	})
}
