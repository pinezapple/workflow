package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"workflow/valkyrie/business"
	"workflow/valkyrie/controller"
	"workflow/valkyrie/core"
	"workflow/valkyrie/model"
	utilsModel "workflow/workflow-utils/model"
)

func HandleGeneratedFile(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &model.HandlerGeneratedFileReq{}, handleGeneratedFile)
}

func ReuploadFileToMinio(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &model.ReuploadToMinioReq{}, reuploadFileToMinio)
}

func handleGeneratedFile(c echo.Context, request interface{}) (statusCode int, data interface{}, err error) {
	//logger init
	logger.Info(" handling generated files")
	req := request.(*model.HandlerGeneratedFileReq)
	mainConf := core.GetMainConfig()
	if mainConf.HardDiskOnly {
		go business.SaveToDBAsync(req)
	} else {
		go business.UploadFileToMinioAsync(req)
	}

	return http.StatusOK, data, nil
}

//
func reuploadFileToMinio(c echo.Context, request interface{}) (statusCode int, data interface{}, err error) {
	//logger init
	logger.Info(" REupload file to minio")
	req := request.(*model.ReuploadToMinioReq)

	userid := c.QueryParam("username")
	ctx := c.Request().Context()

	err = business.ReuploadFileToMinio(ctx, req, userid)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, data, nil
}

// @Summary Get generated files
// @Description Get generated files
// @Produce json
// @Param page_size query int false "Number of result"
// @Param page_token query int false "Current page"
// @Param filter query string false "Filter"
// @Param order query string false "Order"
// @Success 200 {array} model.GetGeneratedFileResponse "List generated files"
// @Tags files
// @Router /files [GET]
func GetGeneratedFiles(c echo.Context) (erro error) {
	//logger init
	logger.Info("get generated files")
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

	rawdata, total, er := business.GetGeneratedFile(ctx, pageSize, pageToken, filter, orderParam)
	if er != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})
	}

	resp := &model.GetGeneratedFileResponse{
		Total: total,
		File:  rawdata,
	}

	//TODO: do some filter with raw data here
	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: resp,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})

}
