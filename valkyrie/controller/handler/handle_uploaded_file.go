package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"workflow/valkyrie/business"
	"workflow/valkyrie/controller"
	"workflow/valkyrie/core"
	"workflow/valkyrie/model"
	utilsModel "workflow/workflow-utils/model"
)

const FILE_UPLOAD_TYPE = "upload"
const FILE_GENERATED_TYPE = "system"

// @Summary Check Resumable chunk is uploaded
// @Description Check resumable chunk is uploaded or not
// @Produce json
// @Param resumableIdentifier query string true "Identifier of upload"
// @Param resumableChunkNumber query string true "Chunk number"
// @Success 200 {string} string "ok"
// @Tags files
// @Router /files/_resumable [GET]
func GetResumableUploadFile(c echo.Context) (err error) {
	r := c.Request()
	config := core.GetMainConfig()
	resumableIdentifier := r.URL.Query()["resumableIdentifier"]
	resumableChunkNumber := r.URL.Query()["resumableChunkNumber"]

	path := fmt.Sprintf("%s/%s", config.UploadTempDirPrefix, resumableIdentifier[0])
	relativeChunk := fmt.Sprintf("%s/%s-%s", path, "part", resumableChunkNumber[0])

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			logger.Errorf("Error when create temporary dir: %s. Reason %s", path, err.Error())
			return
		}
	}

	if _, err = os.Stat(relativeChunk); os.IsNotExist(err) {
		logger.Debugf("Chunk not found: %s", relativeChunk)
		http.Error(c.Response().Writer, http.StatusText(http.StatusNotFound), http.StatusExpectationFailed)
	} else {
		logger.Infof("Chunk already exist: %s", relativeChunk)
		http.Error(c.Response().Writer, "Chunk already exist", http.StatusCreated)
	}
	return nil
}

// @Summary Upload a resumable chunk
// @Description Upload a resumable chunk
// @Produce json
// @Param resumableIdentifier query string true "Identifier of upload"
// @Param resumableChunkNumber query string true "Chunk number"
// @Param resumableFilename query string true "File name"
// @Param resumableTotalChunks query int true "Total chunks of a upload file"
// @Param sample query string false "Sample name"
// @Param workflow query string false "Workflow UUID"
// @Param projectPath query string true "Project Path"
// @Param projectID query string true "Project ID"
// @Success 200 {string} string "ok"
// @Tags files
// @Router /files/_resumable [POST]
func ResumableUploadFile(c echo.Context) (err error) {
	r := c.Request()
	defer r.Body.Close()
	config := core.GetMainConfig()
	userid := r.Context().Value("UserID").(string)
	resumableIdentifier := r.URL.Query()["resumableIdentifier"]
	resumableChunkNumber := r.URL.Query()["resumableChunkNumber"]
	resumableFileName := r.URL.Query()["resumableFilename"][0]

	path := fmt.Sprintf("%s/%s", config.UploadTempDirPrefix, resumableIdentifier[0])
	relativeChunk := fmt.Sprintf("%s/%s-%s", path, "part", resumableChunkNumber[0])

	err = r.ParseMultipartForm(25 << 20)
	if err != nil {
		logger.Errorf("Error when parse multipart form. Reason: %s", err.Error())
		os.RemoveAll(path)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		logger.Errorf("Read from request failed: %s", err.Error())
		os.RemoveAll(path)
		return
	}

	/*
		projectPath := c.QueryParam("projectPath")
		if len(projectPath) == 0 {
			logger.Errorf("No project path is defined")
			os.RemoveAll(path)
			return
		}


		projectID := c.QueryParam("projectID")
		if len(projectID) == 0 {
			logger.Errorf("No project id is defined")
			os.RemoveAll(path)
			return

		}
	*/
	projectPath := ""
	projectID := ""

	defer file.Close()

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			logger.Errorf("Can not create temporary path: %s. Reason: %s", path, err.Error())
			http.Error(c.Response().Writer, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return nil
		}
	}

	f, err := os.OpenFile(relativeChunk, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Errorf("Error when create chunk file: %s. Reason: %s", relativeChunk, err.Error())
		http.Error(c.Response().Writer, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return nil
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		logger.Errorf("Error when copy to temp chunk file: %s", err.Error())
		http.Error(c.Response().Writer, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return nil
	}

	/*
		If it is the last chunk, trigger the recombination of chunks
	*/
	resumableTotalChunks := r.URL.Query()["resumableTotalChunks"]
	current, err := strconv.Atoi(resumableChunkNumber[0])
	total, err := strconv.Atoi(resumableTotalChunks[0])

	if current == total {
		workflowUUID := c.QueryParam("workflow")
		sample := c.QueryParam("sample")
		ctx := c.Request().Context()
		fmt.Println("In current == total")
		/////
		f, err := business.HandleUploadedFile(ctx, userid, path, resumableFileName, sample, workflowUUID, projectPath, projectID)
		if err != nil {
			os.RemoveAll(path)
			return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
				Data: nil,
				Error: utilsModel.ResponseError{
					Message: err.Error(),
					Code:    http.StatusInternalServerError,
				},
			})
		}

		response := &model.GetUserUploadFileList{
			FileUUID: f.FileUUID,
			Path:     f.Path,
			Filename: resumableFileName,
			Filesize: f.Filesize,
		}

		//runtime.GC()
		debug.FreeOSMemory()

		return c.JSON(http.StatusOK, &utilsModel.Response{
			Data: response,
			Error: utilsModel.ResponseError{
				Message: "",
				Code:    http.StatusOK,
			},
		})
	}

	debug.FreeOSMemory()
	//runtime.GC()

	return nil
}

// @Summary Get uploaded and generated files
// @Description Get uploaded and generated files
// @Produce json
// @Param filter query string false "Filter: filter=project_id=4e07845b-5569-4871-b19e-1c8ac90e3052&project_path=/output"
// @Param order query string false "Order"
// @Success 200 {array} model.DataFile "List files"
// @Tags files
// @Router /files [GET]
func GetFiles(c echo.Context) (erro error) {
	logger.Info("Get user files")
	ctx := c.Request().Context()

	// pageSizeString := c.QueryParam("page_size")
	// pageSize, er := strconv.Atoi(pageSizeString)
	// if er != nil {
	// 	return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
	// 		Data: nil,
	// 		Error: utilsModel.ResponseError{
	// 			Message: er.Error(),
	// 			Code:    http.StatusNotAcceptable,
	// 		},
	// 	})
	// }

	// pageTokenString := c.QueryParam("page_token")
	// pageToken, er := strconv.Atoi(pageTokenString)
	// if er != nil {
	// 	return c.JSON(http.StatusNotAcceptable, &utilsModel.Response{
	// 		Data: nil,
	// 		Error: utilsModel.ResponseError{
	// 			Message: er.Error(),
	// 			Code:    http.StatusNotAcceptable,
	// 		},
	// 	})
	// }

	filter := c.QueryParam("filter")
	order := c.QueryParam("order")
	orderParam := strings.ReplaceAll(order, ";", ",")

	rawdata, er := business.GetUserFiles(ctx, filter, orderParam)
	if er != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	}

	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: rawdata,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})
}

// @Summary Get uploaded files
// @Description Get uploaded files
// @Produce json
// @Param page_size query int false "Number of result"
// @Param page_token query int false "Current page"
// @Param filter query string false "Filter"
// @Param order query string false "Order"
// @Success 200 {array} model.GetUploadedFileResponse "List uploaded files"
// @Tags files
// @Router /files [GET]
func GetUserUploadFiles(c echo.Context) (erro error) {
	//logger init
	logger.Info("Get user files")
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

	rawdata, total, er := business.GetUserUploadedFile(ctx, pageSize, pageToken, filter, orderParam)
	if er != nil {
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: er.Error(),
				Code:    http.StatusInternalServerError,
			},
		})
	}
	// TODO: maybe do some filter here with the rawdata
	resp := &model.GetUploadedFileResponse{
		Total: total,
		File:  rawdata,
	}

	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: resp,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})
}

// @Summary Update project path for files
// @Description Update project path for files with uuids
// @Produce json
// @Param updateRequest body model.UpdatePathFiles true "Update project path"
// @Success 200 {string} string "ok"
// @Tags files
// @Router /files/project_path [PUT]
func UpdateFileProjectPath(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &model.UpdatePathFiles{}, updateFileProjectPath)
}

func updateFileProjectPath(c echo.Context, request interface{}) (statusCode int, data interface{}, err error) {
	logger.Info(" update project paht for files ")
	req := request.(*model.UpdatePathFiles)

	if len(req.PathFiles) == 0 {
		return http.StatusOK, nil, nil
	}

	ctx := c.Request().Context()
	err = business.HandleUpdateFileProjectPath(ctx, req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, nil, nil
}
