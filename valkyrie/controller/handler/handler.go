package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"workflow/valkyrie/controller"
	"workflow/valkyrie/core"
	"workflow/valkyrie/dao"
	"workflow/valkyrie/model"
	"workflow/valkyrie/utils"
)

var (
	logger *core.Logger = core.GetLogger()
)

func DeleteFailRunDir(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &model.DeleteFailRunDirReq{}, deleteFailRunDir)
}

func deleteFailRunDir(c echo.Context, request interface{}) (statusCode int, data interface{}, err error) {
	//logger init
	logger.Info(" delete fail run directory")
	req := request.(*model.DeleteFailRunDirReq)

	if req.RunUUID == "" {
		// this is just a delete notification
		if len(req.Filename) == 0 {
			return http.StatusOK, data, nil
		}
		dir := utils.GetTaskDir(req.Filename[0])
		err = utils.DeleteDir(dir)
		if err != nil {
			logger.Error(err.Error())
			return http.StatusInternalServerError, data, err
		}
	} else {
		// this is a first fail-run notification
		mainConf := core.GetMainConfig()
		ctx := c.Request().Context()
		gfDAO := dao.GetGeneratedFileDAO()
		db := core.GetDBObj()

		expiredCloudTime := time.Now().Add(time.Duration(mainConf.NormalFileTTL) * time.Hour)

		err := gfDAO.UpdateCloudExpiredTime(ctx, db, req.UserID, req.RunUUID, expiredCloudTime)
		if err != nil {
			logger.Error(err.Error())
		}

		err = gfDAO.UpdateDoneRun(ctx, db, req.UserID, req.RunUUID)
		if err != nil {
			logger.Error(err.Error())
		}

	}
	return http.StatusOK, data, nil
}

/*
func GetUserFilesWithFilters(c echo.Context) (erro error) {
	//logger init
	logger.Info("Get user files")
	ctx := c.Request().Context()
	db := core.GetDBObj()
	fileDAO := dao.GetFileDAO()
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

	rawdata, er := fileDAO.GetFilesByOffsetWithFilters(ctx, db, userID, pageSize, pageToken, filter, orderParam)
	if er != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	}

	var data []*model.FileHTTPResp
	for i := 0; i < len(rawdata); i++ {
		var path string
		if rawdata[i].UploadSuccess {
			path = utils.ConstructUnsafeDownloadURL(rawdata[i].UserID, rawdata[i].RunUUID, rawdata[i].TaskUUID, rawdata[i].Filename, rawdata[i].Bucket)
		} else {
			path = rawdata[i].LocalPath
		}

		tg := &model.FileHTTPResp{
			Path:      path,
			UserID:    rawdata[i].UserID,
			RunUUID:   rawdata[i].RunUUID,
			RunName:   rawdata[i].RunName,
			TaskUUID:  rawdata[i].TaskUUID,
			TaskName:  rawdata[i].TaskName,
			Filename:  rawdata[i].Filename,
			Filesize:  rawdata[i].Filesize,
			CreatedAt: rawdata[i].CreatedAt,
			ExpiredAt: rawdata[i].ExpiredAt,
		}
		data = append(data, tg)
	}

	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: data,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})

}

type SampleKey struct {
	SampleName string
	CreatedAt  int64
}
type ByCreatedAt []SampleKey

func (samples ByCreatedAt) Len() int           { return len(samples) }
func (samples ByCreatedAt) Swap(i, j int)      { samples[i], samples[j] = samples[j], samples[i] }
func (samples ByCreatedAt) Less(i, j int) bool { return samples[i].CreatedAt > samples[j].CreatedAt }
*/
