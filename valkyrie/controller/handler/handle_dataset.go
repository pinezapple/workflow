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

// @Summary Create dataset
// @Description Create dataset
// @Produce json
// @Param sample body model.CreateNewDatasetReq true "Dataset info"
// @Success 200 {string} string "ok"
// @Tags dataset
// @Router /datasets [POST]
func CreateDataset(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &model.CreateNewDatasetReq{}, createDataset)
}

func createDataset(c echo.Context, request interface{}) (statusCode int, data interface{}, err error) {
	logger.Info(" creating new dataset")
	req := request.(*model.CreateNewDatasetReq)
	ctx := c.Request().Context()

	_, err = business.CreateDataset(ctx, req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, data, nil
}

// @Summary Get datasets
// @Description Get datasets
// @Produce json
// @Param dataset_uuid path string true "dataset UUID"
// @Param page_size query int false "Number of result"
// @Param page_token query int false "Current page"
// @Param filter query string false "Filter"
// @Param order query string false "Order"
// @Success 200 {array} model.Dataset "List datasets" 
// @Tags dataset
// @Router /dataset [GET]
func GetDatasetFilter(c echo.Context) (erro error) {
	logger.Info("Get dataset with filter")
	ctx := c.Request().Context()

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

	dataset, er := business.GetDatasetWithFilter(ctx, pageSize, pageToken, filter, orderParam)
	if er != nil {
		return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: dataset,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})
}
