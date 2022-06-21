package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"workflow/valkyrie/business"
	"workflow/valkyrie/controller"
	"workflow/valkyrie/model"
	utilsModel "workflow/workflow-utils/model"
)

// @Summary Create sample
// @Description Create sample
// @Produce json
// @Param sample body model.CreateNewSampleReq true "Sample info"
// @Success 200 {string} string "ok"
// @Tags samples
// @Router /samples [POST]
func CreateSample(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &model.CreateNewSampleReq{}, createSample)
}

func createSample(c echo.Context, request interface{}) (statusCode int, data interface{}, err error) {
	logger.Info(" creating new sample")
	req := request.(*model.CreateNewSampleReq)
	ctx := c.Request().Context()

	err = business.CreateSample(ctx, req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, data, nil
}

// @Summary Get samples by dataset
// @Description Get samples by dataset
// @Produce json
// @Param dataset_uuid path string true "dataset UUID"
// @Param page_size query int false "Number of result"
// @Param page_token query int false "Current page"
// @Param filter query string false "Filter"
// @Param order query string false "Order"
// @Success 200 {array} model.SampleDetailResp "List sample files"
// @Tags samples
// @Router /samples/workflow/:dataset_uuid [GET]
func GetSampleByDataset(c echo.Context) (erro error) {
	logger.Info("Get samples by dataset with filter")
	ctx := c.Request().Context()
	//userID := ctx.Value("UserID").(string)
	// datasetUUID := c.QueryParam("dataset_uuid")
	datasetUUID := c.Param("dataset_uuid")

	pageSizeString := c.QueryParam("page_size")
	pageSize, er := strconv.Atoi(pageSizeString)
	if er != nil {
		return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	pageTokenString := c.QueryParam("page_token")
	pageToken, er := strconv.Atoi(pageTokenString)
	if er != nil {
		return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	filter := c.QueryParam("filter")
	order := c.QueryParam("order")
	orderParam := strings.ReplaceAll(order, ";", ",")

	sample, err := business.GetSampleByDataset(ctx, datasetUUID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	}
	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: sample,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})
}

// @Summary Get samples by workflow
// @Description Get samples by workflow
// @Produce json
// @Param workflow_uuid path string true "Workflow UUID"
// @Param page_size query int false "Number of result"
// @Param page_token query int false "Current page"
// @Param filter query string false "Filter"
// @Param order query string false "Order"
// @Success 200 {array} model.GetSampleResponse "List sample files"
// @Tags samples
// @Router /samples/workflow/:workflow_uuid [GET]
func GetSamplesByWorkflow(c echo.Context) (erro error) {
	logger.Info("Get samples by workflow id with filter")
	ctx := c.Request().Context()
	userID := ctx.Value("UserID").(string)
	// workflowUUID := c.QueryParam("workflow_uuid")
	workflowUUID := c.Param("workflow_uuid")

	pageSizeString := c.QueryParam("page_size")
	pageSize, er := strconv.Atoi(pageSizeString)
	if er != nil {
		return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	pageTokenString := c.QueryParam("page_token")
	pageToken, er := strconv.Atoi(pageTokenString)
	if er != nil {
		return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	filter := c.QueryParam("filter")
	order := c.QueryParam("order")
	orderParam := strings.ReplaceAll(order, ";", ",")

	sample, total, err := business.GetSampleByWorkflow(ctx, userID, workflowUUID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	}
	resp := &model.GetSampleResponse{
		Total:   total,
		Samples: sample,
	}

	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: resp,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})
}

// @Summary Get sample
// @Description Get sample
// @Produce json
// @Param page_size query int false "Number of result"
// @Param page_token query int false "Current page"
// @Param filter query string false "Filter"
// @Param order query string false "Order"
// @Success 200 {array} model.SampleDetailResp "List sample files"
// @Tags samples
// @Router /samples [GET]
func GetSampleFilter(c echo.Context) (erro error) {
	logger.Info("Get samples by workflow id with filter")
	ctx := c.Request().Context()
	userID := ctx.Value("UserID").(string)

	pageSizeString := c.QueryParam("page_size")
	pageSize, er := strconv.Atoi(pageSizeString)
	if er != nil {
		return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	pageTokenString := c.QueryParam("page_token")
	pageToken, er := strconv.Atoi(pageTokenString)
	if er != nil {
		return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	filter := c.QueryParam("filter")
	order := c.QueryParam("order")
	orderParam := strings.ReplaceAll(order, ";", ",")

	sample, err := business.GetSampleFilter(ctx, userID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	}

	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: sample,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})
}

// @Summary Get sample detail
// @Description Get sample detail
// @Produce json
// @Param sample_uuid path string true "Sample uuid"
// @Success 200 {} model.SampleDetailResp "List sample files"
// @Tags samples
// @Router /samples/:sample_uuid [GET]
func GetSampleDetail(c echo.Context) (erro error) {
	logger.Info("Get samples by workflow id with filter")
	ctx := c.Request().Context()
	//userID := ctx.Value("UserID").(string)
	// sampleUUID := c.QueryParam("sample_id")
	sampleUUID := c.Param("sample_uuid")
	resp, err := business.GetSampleDetail(ctx, sampleUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	}

	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: resp,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})
}

// @Summary Check sample name is existed
// @Description Check sample name is existed
// @Produce json
// @Param sample_name path string true "Sample name"
// @Param workflow query string true "Workflow UUID"
// @Success 200 {} nil "List sample files"
// @Tags samples
// @Router /samples/name/:sample_name [GET]
func CheckExistSample(c echo.Context) error {
	sampleName := c.Param("sample_name")
	workflow := c.QueryParam("workflow")
	ctx := c.Request().Context()
	userid := ctx.Value("UserID").(string)

	ok, err := business.CheckExistSample(ctx, sampleName, workflow, userid)
	if err != nil {
		return c.JSON(http.StatusBadGateway, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	} else {
		if ok {
			return c.JSON(http.StatusBadGateway, &utilsModel.Response{
				Data: nil,
				Error: utilsModel.ResponseError{
					Message: "Sample name existed",
					Code:    http.StatusInternalServerError,
				},
			})
		}
		return c.JSON(http.StatusOK, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: "",
				Code:    http.StatusOK,
			},
		})
	}
}
