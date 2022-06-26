package handler

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
	"workflow/valkyrie/controller"
	"workflow/valkyrie/core"
	"workflow/valkyrie/dao"
	"workflow/valkyrie/model"
	utilsModel "workflow/workflow-utils/model"
)

// @Summary Download file
// @Description Download file
// @Param file_name path string true "File name"
// @Param run_uuid query string true "Run UUID"
// @Param task_uuid query string true "Task UUID"
// @Tags files
// @Router /files/:file_name [GET]
func DownloadFileFromHardDisk(c echo.Context) (erro error) {
	//logger init
	logger.Info("Download user files")

	ctx := c.Request().Context()
	db := core.GetDBObj()
	gfDAO := dao.GetGeneratedFileDAO()
	userID := "tungnt99"

	runuuid := c.QueryParam("run_uuid")
	taskuuid := c.QueryParam("task_uuid")
	// filename := c.QueryParam("file_name")
	filename := c.Param("file_name")

	core.DownloadLock.RLock()
	defer core.DownloadLock.RUnlock()
	file, err := gfDAO.GetFileByFilename(ctx, db, userID, runuuid, taskuuid, filename)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})

	}

	return c.Attachment(file.Path, filename)
}

func GetPresignedDownloadURL(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &model.CreatePresignedURLReq{}, getPresignedDownloadURL)
}

func getPresignedDownloadURL(c echo.Context, request interface{}) (statusCode int, data interface{}, err error) {
	//logger init
	logger.Info("get presigned download url")
	req := request.(*model.CreatePresignedURLReq)
	ctx := c.Request().Context()
	userid := "tungnt99"

	// Get minio client
	s3Client := core.GetMinIOClient()
	db := core.GetDBObj()
	gfDAO := dao.GetGeneratedFileDAO()

	// Get bucket from db
	bucket, err := gfDAO.GetFileBucket(ctx, db, req.FileUUID)
	if err != nil {
		logger.Errorf("Get file bucket error: %s", err.Error())
		return http.StatusInternalServerError, nil, err
	}
	var key string
	if req.RunUUID == "" && req.TaskUUID == "" {
		key = userid + "/" + req.Filename
	} else {
		key = userid + "/" + req.RunUUID + "/" + req.TaskUUID + "/" + req.Filename
	}

	// Start Geting from s3
	url, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := url.Presign(time.Duration(req.TTL) * time.Minute)
	if err != nil {
		logger.Errorf("Get presign url error: %s", err.Error())
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, urlStr, nil
}

// @Summary Delete file
// @Description Delete file
// @Param file_id path string true "File id"
// @Tags files
// @Router /files/:file_id [DELETE]
func DeleteFileFromHardDisk(c echo.Context) (erro error) {
	ctx := c.Request().Context()
	db := core.GetDBObj()
	fileDAO := dao.GetFileDAO()
	userID := "tungnt99"

	fileId := c.Param("file_id")
	logger.Infof("Delete file: %s", fileId)
	logger.Infof("Path params names: %v", c.ParamNames())

	err := fileDAO.DeleteFileFromHardDisk(ctx, db, userID, fileId, true)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusNotAcceptable,
			},
		})

	}

	return c.JSON(http.StatusOK, nil)
}
