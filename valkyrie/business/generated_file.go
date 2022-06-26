package business

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"workflow/valkyrie/api/heimdall"
	"workflow/valkyrie/core"
	"workflow/valkyrie/dao"
	"workflow/valkyrie/model"
	"workflow/valkyrie/utils"
	"gorm.io/gorm"
)

func SaveToDBAsync(req *model.HandlerGeneratedFileReq) {
	ctx := context.Background()
	// Get minio client
	db := core.GetDBObj()
	mainConf := core.GetMainConfig()
	gfDAO := dao.GetGeneratedFileDAO()

	fileLastName := core.GetFileName(req.Filename)
	if len(req.Filename) == 0 {
		logger.Error(" no generated files to handle")
		return
	}

	for i := 0; i < len(req.Filename); i++ {
		f := &model.GeneratedFile{
			FileUUID:      uuid.New().String(),
			UserID:        req.UserID,
			RunID:         req.RunID,
			RunUUID:       req.RunUUID,
			RunName:       req.RunName,
			ProjectID:     req.ProjectID,
			ProjectPath:   req.ProjectPath,
			TaskID:        req.Taskid,
			TaskUUID:      req.TaskUUID,
			TaskName:      req.TaskName,
			Path:          req.Filename[i],
			Filename:      fileLastName[i],
			Filesize:      req.Filesize[i],
			UploadSuccess: false,
			DoneRun:       false,
			CreatedAt:     time.Now(),
		}

		err := gfDAO.SaveFile(ctx, db, f)
		if err != nil {
			logger.Errorf("Save generated file info to db err : %s", err.Error())
		}
	}

	if req.LastTask {
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

	return
}

func UploadFileToMinioAsync(req *model.HandlerGeneratedFileReq) {
	ctx := context.Background()
	mainConf := core.GetMainConfig()

	// Get minio client
	s3Client := core.GetMinIOClient()
	db := core.GetDBObj()
	bucketDAO := dao.GetBucketDAO()
	gfDAO := dao.GetGeneratedFileDAO()

	fileLastName := core.GetFileName(req.Filename)
	if len(req.Filename) == 0 {
		logger.Error(" no files to upload")
		return
	}

	var success, fail []*model.FileUploadResp

	for i := 0; i < len(req.Filename); i++ {
		f := &model.GeneratedFile{
			FileUUID:      uuid.New().String(),
			UserID:        req.UserID,
			RunID:         req.RunID,
			RunUUID:       req.RunUUID,
			RunName:       req.RunName,
			ProjectID:     req.ProjectID,
			ProjectPath:   req.ProjectPath,
			TaskID:        req.Taskid,
			TaskUUID:      req.TaskUUID,
			TaskName:      req.TaskName,
			Path:          req.Filename[i],
			Filename:      fileLastName[i],
			Filesize:      req.Filesize[i],
			UploadSuccess: false,
			DoneRun:       false,
			CreatedAt:     time.Now(),
		}

		fresp := &model.FileUploadResp{
			Name: fileLastName[i],
			Type: "FILE",
		}

		key := req.UserID + "/" + req.RunUUID + "/" + req.TaskUUID + "/" + fileLastName[i]

		bucket, iter, newBucket := core.GetMinioBucket(req.Filesize[i])
		if newBucket {
			err := core.MountBucketToDir(bucket, mainConf.FUSEMountpoint+"/"+bucket, mainConf.MinioAuthenFile, mainConf.MinioEndpoint)
			if err != nil {
				logger.Errorf("Mount bucket to directory error: %s", err.Error())
				core.DeleteFromBucket(iter, req.Filesize[i])
				fail = append(fail, fresp)

				handleMinioUploadErr(ctx, db, gfDAO, f)
			}
			continue
		}

		// Open the file for use
		file, err := os.Open(req.Filename[i])
		if err != nil {
			logger.Errorf("Open file error: %s", err.Error())
			core.DeleteFromBucket(iter, req.Filesize[i])
			fail = append(fail, fresp)

			handleMinioUploadErr(ctx, db, gfDAO, f)
			file.Close()
			continue
		}
		defer file.Close()

		// Put object to s3
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Body:   file,
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			logger.Errorf("Upload file to S3 error: %s", err.Error())
			core.DeleteFromBucket(iter, req.Filesize[i])
			fail = append(fail, fresp)

			handleMinioUploadErr(ctx, db, gfDAO, f)
			continue

		}

		// TODO: save bucket and file to db
		if newBucket {
			bucketDAO.CreateNewBucket(ctx, db, bucket, req.Filesize[i], iter)
		} else {
			bucketDAO.AddToBucket(ctx, db, bucket, req.Filesize[i])
		}

		f.UploadSuccess = true
		err = gfDAO.SaveFile(ctx, db, f)
		if err != nil {
			logger.Errorf("Save upload file error: %s", err.Error())
			core.DeleteFromBucket(iter, req.Filesize[i])
			fail = append(fail, fresp)
			continue
		}

		fresp.Path = mainConf.MinioEndpoint + "/" + bucket + "/" + key

		success = append(success, fresp)
	}

	err := heimdall.UploadFileDoneStatus(req.TaskUUID, success, fail)
	if err != nil {
		logger.Errorf("Send file info to heimdall error: %s", err.Error())
	}

	if req.LastTask {
		handleUploadFileOfLastTask(ctx, db, gfDAO, req)
	}
}

func handleMinioUploadErr(ctx context.Context, db *gorm.DB, gfDAO dao.IgenFileDAO, f *model.GeneratedFile) {
	f.UploadSuccess = false
	// f.UploadExpiredAt = time.Now().Add(2 * 24 * time.Hour)
	err := gfDAO.SaveFile(ctx, db, f)
	if err != nil {
		logger.Errorf("Save upload file error: %s", err.Error())
	}
}

func handleUploadFileOfLastTask(ctx context.Context, db *gorm.DB, gfDAO dao.IgenFileDAO, req *model.HandlerGeneratedFileReq) {
	mainConf := core.GetMainConfig()
	expiredCloudTime := time.Now().Add(time.Duration(mainConf.NormalFileTTL) * time.Hour)
	expiredUploadTime := time.Now().Add(time.Duration(mainConf.UploadPenalty) * time.Hour)

	// update cloud expired time for succesfully uploaded file
	err := gfDAO.UpdateCloudExpiredTime(ctx, db, req.UserID, req.RunUUID, expiredCloudTime)
	if err != nil {
		logger.Error(err.Error())
	}

	// update that this run is done
	err = gfDAO.UpdateDoneRun(ctx, db, req.UserID, req.RunUUID)
	if err != nil {
		logger.Error(err.Error())
	}

	filename := req.Filename[0]

	if !mainConf.HardDiskOnly {
		fileNotUploaded, err := gfDAO.GetFileNotUploadedByRunID(ctx, db, req.UserID, req.RunUUID)
		if err != nil {
			logger.Error(err.Error())
		}

		// if all file are uploaded, delete file from hard disk
		if len(fileNotUploaded) == 0 {
			reg := utils.GetTaskDirectoryRegexp(req.Taskid, filename)
			err := utils.DeleteDirWithRegex(reg)
			if err != nil {
				logger.Error(err.Error())
			}

		} else {
			// update hard disk expired time for file that is not uploaded
			err = gfDAO.UpdateHardDiskExpiredTime(ctx, db, req.UserID, req.RunUUID, expiredUploadTime)
			if err != nil {
				logger.Error(err.Error())
			}

			// delete files that are uploaded
			fileToDelete, err := gfDAO.GetFileUploadedByRunID(ctx, db, req.UserID, req.RunUUID)
			if err != nil {
				logger.Error(err.Error())
			}
			for i := 0; i < len(fileToDelete); i++ {
				err := utils.DeleteFile(fileToDelete[i].Path)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}
	}
}

func ReuploadFileToMinio(ctx context.Context, req *model.ReuploadToMinioReq, userid string) (err error) {
	mainConf := core.GetMainConfig()
	// Get minio client
	s3Client := core.GetMinIOClient()
	db := core.GetDBObj()
	//fileDAO := dao.GetFileDAO()
	gfDAO := dao.GetGeneratedFileDAO()

	bucketDAO := dao.GetBucketDAO()

	//TODO: Change this to file UUID
	tobeUploadedFile, err := gfDAO.GetFileByFilename(ctx, db, userid, req.RunUUID, req.TaskUUID, req.Filename)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	bucket, iter, newBucket := core.GetMinioBucket(tobeUploadedFile.Filesize)
	if newBucket {
		err := core.MountBucketToDir(bucket, mainConf.FUSEMountpoint+"/"+bucket, mainConf.MinioAuthenFile, mainConf.MinioEndpoint)
		if err != nil {
			logger.Error(err.Error())
			core.DeleteFromBucket(iter, tobeUploadedFile.Filesize)
			return err
		}
	}

	key := req.UserID + "/" + req.RunUUID + "/" + req.TaskUUID + "/" + req.Filename

	// Open the file for use
	file, err := os.Open(tobeUploadedFile.Path)
	if err != nil {
		logger.Error(err.Error())
		core.DeleteFromBucket(iter, tobeUploadedFile.Filesize)
		return err
	}
	defer file.Close()

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		core.DeleteFromBucket(iter, tobeUploadedFile.Filesize)
		return err
	}

	// TODO: save bucket and file to db
	if newBucket {
		bucketDAO.CreateNewBucket(ctx, db, bucket, tobeUploadedFile.Filesize, iter)
	} else {
		bucketDAO.AddToBucket(ctx, db, bucket, tobeUploadedFile.Filesize)
	}

	// update as uploaded
	err = gfDAO.UpdateFileAsUploaded(ctx, db, userid, tobeUploadedFile.FileUUID, bucket)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if tobeUploadedFile.DoneRun {
		// if all uploaded then delete all, else only delete file
		fileNotUploaded, err := gfDAO.GetFileNotUploadedByRunID(ctx, db, req.UserID, req.RunUUID)
		if err != nil {
			logger.Error(err.Error())
		}

		if len(fileNotUploaded) == 0 {
			reg := utils.GetTaskDirectoryRegexp(tobeUploadedFile.TaskID, tobeUploadedFile.Path)
			err := utils.DeleteDirWithRegex(reg)
			if err != nil {
				logger.Error(err.Error())
			}

		} else {
			err = utils.DeleteFile(tobeUploadedFile.Path)
			if err != nil {
				logger.Error(err.Error())
			}
		}

	}

	return
}

func GetGeneratedFile(ctx context.Context, pageSize, pageToken int, filter, orderParam string) (data []model.GeneratedFile, total int64, err error) {
	db := core.GetDBObj()
	fileDAO := dao.GetGeneratedFileDAO()
	userID := "tungnt99"

	rawdata, total, er := fileDAO.GetUserGeneratedFiles(ctx, db, userID, pageSize, pageToken, filter, orderParam)
	if er != nil {
		logger.Error(er.Error())
		return nil, 0, er
	}

	return rawdata, total, nil
}
