package business

import (
	"context"
	"strconv"
	"time"

	"github.com/vfluxus/valkyrie/core"
	"github.com/vfluxus/valkyrie/dao"
	"github.com/vfluxus/valkyrie/utils"
	"github.com/vfluxus/workflow-utils/model"
)

var (
	logger *core.Logger = core.GetLogger()
)

func sleepContext(ctx context.Context, delay time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}

func DailyCheckOnExpiredObjOnHardDisk(parentCtx context.Context) (fn model.Daemon, err error) {
	fn = func() {

		logger.Info("Daily check for expired obj on hardisk daemon starting")

		mainConf := core.GetMainConfig()
		db := core.GetDBObj()
		fileDAO := dao.GetGeneratedFileDAO()

		for {
			sleepContext(parentCtx, time.Duration(mainConf.UpdateStorageChangeInterval)*time.Hour)

			core.DownloadLock.Lock()
			if parentCtx.Err() != nil {
				core.DownloadLock.Unlock()
				return
			}

			// Get expiried list from db
			files, err := fileDAO.GetExpiredFiles(parentCtx, db, time.Now())
			if err != nil {
				logger.Errorf("Get expired files error: %s", err.Error())
			}

			for i := 0; i < len(files); i++ {
				// Start delete sequence
				err := fileDAO.DeleteFile(parentCtx, db, files[i].UserID, files[i].RunUUID, files[i].TaskUUID, files[i].Filename)
				if err != nil {
					logger.Errorf("Delete file error: %s", err.Error())
				}

				if files[i].RunUUID != "" {
					fileNotUploaded, err := fileDAO.GetExpiredFilesByRunID(parentCtx, db, files[i].UserID, files[i].RunUUID, time.Now())
					if err != nil {
						logger.Errorf("Get expired files by run id error: %s", err.Error())
					}

					if len(fileNotUploaded) == 0 {
						/*
							err := utils.DeleteDir(utils.GetParentDir(files[i].RunUUID, files[i].UserID))
							if err != nil {
								lg.LogErr(err)
							}err := utils.DeleteDir(utils.GetParentDir(files[i].RunUUID, files[i].UserID))
							if err != nil {
								lg.LogErr(err)
							}
						*/
						reg := utils.GetTaskDirectoryRegexp(files[i].TaskID, files[i].Path)
						err := utils.DeleteDirWithRegex(reg)
						if err != nil {
							logger.Errorf("Delete dir with regex error: %s", err.Error())
						}
					}
				}

				err = utils.DeleteFile(files[i].Path)
				if err != nil {
					logger.Errorf("Delete file error: %s", err.Error())
				}
			}
			core.DownloadLock.Unlock()
		}
	}
	return fn, nil
}

/*
// TODO: fix this later
func DailyCheckOnExpiredObj(parentCtx context.Context) (fn model.Daemon, err error) {
	fn = func() {

		logger.Info("Daily check for expired obj daemon starting")

		mainConf := core.GetMainConfig()
		db := core.GetDBObj()
		fileDAO := dao.GetFileDAO()
		bucketDAO := dao.GetBucketDAO()

		for {
			sleepContext(parentCtx, time.Duration(mainConf.UpdateStorageChangeInterval)*time.Hour)
			core.DownloadLock.Lock()
			if parentCtx.Err() != nil {
				core.DownloadLock.Unlock()
				return
			}

			// Get expiried list from db
			files, err := fileDAO.GetExpiredFiles(parentCtx, db, time.Now())
			if err != nil {
				logger.Errorf("Get expired files error: %s", err.Error())
			}

			for i := 0; i < len(files); i++ {
				// Start delete sequence
				err := fileDAO.DeleteFile(parentCtx, db, files[i].UserID, files[i].RunUUID, files[i].TaskUUID, files[i].Filename, files[i].Safe)
				if err != nil {
					logger.Errorf("Delete file error: %s", err.Error())
				}

				// start to get update for buckets
				err = bucketDAO.DeleteFromBucket(parentCtx, db, files[i].Bucket, files[i].Filesize)
				if err != nil {
					logger.Errorf("Delete from bucket error: %s", err.Error())
				}

				if !files[i].Safe {
					_, err = DeleteFromStorage(files[i].Bucket, files[i].UserID+"/"+files[i].RunUUID+"/"+files[i].TaskUUID+"/"+files[i].Filename, files[i].Filesize)
					if err != nil {
						logger.Errorf("Delete from storage error: %s", err.Error())
					}
				}
			}
			core.DownloadLock.Unlock()
		}
	}
	return fn, nil
}
*/

/*
// TODO: fix this later
func CheckOnExpiredUploadObj(parentCtx context.Context) (fn model.Daemon, err error) {
	fn = func() {

		logger.Info("Daily check for expired obj daemon on")

		mainConf := core.GetMainConfig()
		db := core.GetDBObj()
		fileDAO := dao.GetFileDAO()

		for {
			sleepContext(parentCtx, time.Duration(mainConf.ExpiredUpdateInterval)*time.Hour)

			core.DownloadLock.Lock()
			if parentCtx.Err() != nil {
				core.DownloadLock.Unlock()
				return
			}

			files, err := fileDAO.GetExpiredUploadFiles(parentCtx, db, time.Now())
			if err != nil {
				logger.Errorf("Get expired upload files error: %s", err.Error())
			}

			for i := 0; i < len(files); i++ {
				// Start delete sequence
				err := fileDAO.DeleteFile(parentCtx, db, files[i].UserID, files[i].RunUUID, files[i].TaskUUID, files[i].Filename, files[i].Safe)
				if err != nil {
					logger.Errorf("Delete file: %s", err.Error())
				}

				if files[i].RunUUID != "" {
					fileNotUploaded, err := fileDAO.GetFileNotUploadedByRunID(parentCtx, db, files[i].UserID, files[i].RunUUID)
					if err != nil {
						logger.Errorf("Get file not uploaed by run id error: %s", err.Error())
					}

					if len(fileNotUploaded) == 0 {
						reg := utils.GetTaskDirectoryRegexp(files[i].TaskID, files[i].LocalPath)
						err := utils.DeleteDirWithRegex(reg)
						if err != nil {
							logger.Errorf("Delete dir with regex error: %s", err.Error())
						}
					}
				}

				err = utils.DeleteFile(files[i].LocalPath)
				if err != nil {
					logger.Errorf("Delete file error: %s", err.Error())
				}

			}
			core.DownloadLock.Unlock()
		}
	}
	return fn, nil
}
*/
func RestartDaemon(parentCtx context.Context) {
	logger.Info("Reloading old data")

	mainConf := core.GetMainConfig()
	// start restart sequence
	// read buckets state from db
	db := core.GetDBObj()
	bucketDAO := dao.GetBucketDAO()
	buckets, err := bucketDAO.SelectAllBucket(parentCtx, db)
	if err != nil {
		panic(err)
	}
	var size []int64
	var count []int32
	for i := 0; i < len(buckets); i++ {
		size = append(size, buckets[i].CurrentSize)
		count = append(count, buckets[i].CurrentCount)
		err := core.MountBucketToDir("bucket"+strconv.Itoa(i), mainConf.FUSEMountpoint+"/bucket"+strconv.Itoa(i), mainConf.MinioAuthenFile, mainConf.MinioEndpoint)
		if err != nil {
			panic(err)
		}
	}

	core.RestartMinioBuckets(size, count)
	/*
		fn = func() {
			<-parentCtx.Done()

			lg.LogInfo("Shutting down service")
		}
		return fn, nil
	*/
}
