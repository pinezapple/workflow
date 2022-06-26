package business

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"workflow/valkyrie/core"
	"workflow/valkyrie/dao"
	"workflow/valkyrie/model"
)

func HandleUploadedFile(ctx context.Context, userid, path, resumableFileName, sample, workflowUUID, projectPath, projectID string) (f *model.UploadedFile, err error) {
	logger.Infof("Combining chunks into one file")
	filePath, err := combineMultipeChunks(path, userid, resumableFileName)
	if err != nil {
		return nil, err
	}

	fileInfo, _ := os.Lstat(filePath)
	db := core.GetDBObj()

	f = &model.UploadedFile{
		FileUUID:    uuid.New().String(),
		UserID:      userid,
		Path:        filePath,
		Filename:    resumableFileName,
		Filesize:    fileInfo.Size(),
		ProjectPath: projectPath,
		ProjectID:   projectID,
		Safe:        true,
		Deleted:     false,
		CreatedAt:   time.Now(),
	}

	// save file
	ufDAO := dao.GetUploadedFileDAO()
	//err = fileDAO.SaveFile(c.Request().Context(), db, f)
	_, err = ufDAO.SaveFile(ctx, db, f)
	if err != nil {
		logger.Errorf("save uploaded file to db error: %s", err.Error())
		return nil, err
	}

	// if upload with sample name
	if sample != "" {
		sampleDAO := dao.GetSampleDAO()
		// check if this sample exist
		s, err := sampleDAO.GetSampleMetadataBySampleName(ctx, db, sample, userid, workflowUUID)
		if err != nil {
			logger.Errorf("save sample to db error: %s", err.Error())
			return nil, err
		}

		fmt.Printf("Handling uploading new sample %v \n", sample)
		// if exist then add to sample
		if s.SampleUUID != "" {
			content, err := sampleDAO.GetSampleFileIndex(ctx, db, s.SampleUUID)
			if err != nil {
				logger.Errorf("save sample to db error: %s", err.Error())
				return nil, err
			}

			err = sampleDAO.AddFileToExistedSample(ctx, db, s, f.FileUUID, content[0].SampleIndex+1)
			if err != nil {
				logger.Errorf("save sample to db error: %s", err.Error())
				return nil, err
			}
		} else {
			// if not exist, create new one
			newSample := &model.Sample{
				SampleUUID:   uuid.New().String(),
				UserID:       userid,
				SampleName:   sample,
				WorkflowUUID: workflowUUID,
				CreatedAt:    time.Now(),
			}

			var files []string
			files = append(files, f.FileUUID)
			_, err = sampleDAO.SaveSample(ctx, db, newSample, files)
			if err != nil {
				logger.Errorf("save sample to db error: %s", err.Error())
				return nil, err
			}

		}
	}
	return f, nil
}

func combineMultipeChunks(chunksDir string, userid string, fileName string) (filePath string, err error) {
	now := time.Now().Unix()
	mainConf := core.GetMainConfig()

	finalPath := mainConf.InputDirPrefix + "/" + userid + "/" + strconv.FormatInt(now, 10) + "/" + fileName

	_ = os.MkdirAll(mainConf.InputDirPrefix+"/"+userid+"/"+strconv.FormatInt(now, 10), os.ModePerm)
	f, err := os.Create(finalPath)
	if err != nil {
		logger.Errorf("Create combine file error : %s. Reason %s", finalPath, err.Error())
		return
	}
	defer f.Close()

	var writeOffset int64 = 0
	files, err := ioutil.ReadDir(chunksDir)

	defer os.RemoveAll(chunksDir)

	for i := 1; i <= len(files); i++ {
		chunkPath := fmt.Sprintf("%s/%s-%d", chunksDir, "part", i)
		logger.Debugf("Merge path: %s", chunkPath)

		if _, err = os.Stat(chunkPath); os.IsNotExist(err) {
			logger.Errorf("Chunk path not found: %s", chunkPath)
			continue
		}

		chunkData, er := ioutil.ReadFile(chunkPath)
		if er != nil {
			logger.Errorf("Error when read chunk part: %s", chunkPath)
			return "", er
		}

		size, er := f.WriteAt(chunkData, writeOffset)
		if er != nil {
			logger.Errorf("Error when write chunk part: %s", chunkPath)
			return "", er
		}
		if size != len(chunkData) {
			logger.Errorf("Number of written bytes differences: %d %d", size, len(chunkData))
		}

		writeOffset = writeOffset + int64(size)
	}

	logger.Infof("Complete combine chunks of file: %s", finalPath)
	return finalPath, nil
}

func GetUserFiles(ctx context.Context, filter, orderParam string) (data []model.DataFile, err error) {
	db := core.GetDBObj()
	fileDAO := dao.GetUploadedFileDAO()
	userID := "tungnt99"

	rawdata, err := fileDAO.GetUserFiles(ctx, db, userID, filter, orderParam)
	if err != nil {
		logger.Errorf("get user files error: %s", err.Error())
		return nil, err
	}

	return rawdata, nil
}

func GetUserUploadedFile(ctx context.Context, pageSize, pageToken int, filter, orderParam string) (data []model.UploadedFile, total int64, err error) {
	db := core.GetDBObj()
	fileDAO := dao.GetUploadedFileDAO()
	userID := "tungnt99"

	rawdata, total, er := fileDAO.GetUserUploadFiles(ctx, db, userID, pageSize, pageToken, filter, orderParam)
	if er != nil {
		logger.Error(err.Error())
		return nil, 0, err
	}
	return rawdata, total, nil
}

func HandleUpdateFileProjectPath(ctx context.Context, pathFiles *model.UpdatePathFiles) (err error) {
	db := core.GetDBObj()
	fileDAO := dao.GetUploadedFileDAO()
	userID := "tungnt99"

	err = fileDAO.UpdateFileProjectPath(ctx, db, pathFiles, userID)
	if err != nil {
		logger.Error(err.Error())
	}

	return err
}
