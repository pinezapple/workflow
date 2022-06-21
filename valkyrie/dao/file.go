package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"workflow/valkyrie/core"
	"workflow/valkyrie/model"
	"workflow/valkyrie/utils"
)

type fileDAO struct{}

type IfileDAO interface {
	SaveFile(ctx context.Context, db *gorm.DB, file *model.File) (err error)
	DeleteFile(ctx context.Context, db *gorm.DB, userid, runid, taskid, filename string, save bool) (err error)
	UpdateCloudExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error)
	UpdateHardDiskExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error)
	UpdateDoneRun(ctx context.Context, db *gorm.DB, userid, runid string) (err error)

	GetFilesByOffSet(ctx context.Context, db *gorm.DB, userid string, offset, limit int64) (result []model.File, err error)
	GetFilesWithClue(ctx context.Context, db *gorm.DB, userid string, clue string) (result []model.File, err error)

	GetFileBucket(ctx context.Context, db *gorm.DB, userid, runid, taskid, filename string) (result string, err error)
	GetExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.File, err error)
	GetExpiredFilesByRunID(ctx context.Context, db *gorm.DB, userid, runid string, date time.Time) (result []model.File, err error)
	GetFilesByOffsetWithFilters(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.File, err error)

	GetFilesByOffsetWithFiltersForUserUploadFile(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.File, total int64, err error)
	GetUserUploadFiles(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string, workflow_uuid string) (result []model.UserUploadFile, err error)

	GetFilesByOffsetWithFiltersForSystemGenFile(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.File, err error)

	GetFileNotUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.File, err error)
	GetFileUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.File, err error)

	GetFileByFilename(ctx context.Context, db *gorm.DB, userid string, runid string, taskid string, filename string) (result *model.File, err error)
	GetExpiredUploadFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.File, err error)
	UpdateFileAsUploaded(ctx context.Context, db *gorm.DB, userid, runid, taskid, bucket, filename string) (err error)

	SampleExisted(ctx context.Context, db *gorm.DB, userid, sample_name, workflow string) (bool, error)

	DeleteFileFromHardDisk(ctx context.Context, db *gorm.DB, userid, fileId string, removeFile bool) (err error)
}

var (
	iFileDAO IfileDAO = &fileDAO{}
)

func SetFileDAO(v IfileDAO) {
	iFileDAO = v
	return
}

func GetFileDAO() (dao IfileDAO) {
	return iFileDAO
}

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- IMPLEMENT ----------------------------------------------------------

func (f *fileDAO) SaveFile(ctx context.Context, db *gorm.DB, file *model.File) (err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Create(file).Error; err != nil {
		return err
	}

	return nil
}

func (f *fileDAO) DeleteFile(ctx context.Context, db *gorm.DB, userid, runid, taskid, filename string, save bool) (err error) {
	// update field
	if save {
		var (
			deleteFile = &model.File{
				UserID:   userid,
				RunUUID:  runid,
				TaskUUID: taskid,
				Filename: filename,
				Deleted:  true,
			}
		)

		if err := db.Model(&model.File{}).WithContext(ctx).Updates(deleteFile).Error; err != nil {
			return err
		}
	}

	// if not save, then run delete query
	if err := db.Model(&model.File{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND task_uuid = ? AND filename = ?", userid, runid, taskid, filename).Delete(&model.Bucket{}).Error; err != nil {
		return err
	}

	return nil
}

func (f *fileDAO) UpdateCloudExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error) {
	mainConf := core.GetMainConfig()
	file := &model.File{
		ExpiredAt: expiredAt,
	}

	if err := db.Model(&model.File{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND upload_success = ?", userid, runid, !mainConf.HardDiskOnly).Updates(file).Error; err != nil {
		return err
	}

	return nil
}

func (f *fileDAO) UpdateHardDiskExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error) {
	file := &model.File{
		UploadExpiredAt: expiredAt,
	}

	if err := db.Model(&model.File{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND upload_success = ?", userid, runid, false).Updates(file).Error; err != nil {
		return err
	}

	return nil
}

func (f *fileDAO) UpdateDoneRun(ctx context.Context, db *gorm.DB, userid, runid string) (err error) {
	file := &model.File{
		DoneRun: true,
	}

	if err := db.Model(&model.File{}).WithContext(ctx).Where("username = ? AND run_uuid = ?", userid, runid).Updates(file).Error; err != nil {
		return err
	}

	return nil
}

func (f *fileDAO) GetFilesByOffSet(ctx context.Context, db *gorm.DB, userid string, offset, limit int64) (result []model.File, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Limit(int(limit)).Offset(int(offset)).Where("username = ?", userid).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (f *fileDAO) GetFilesWithClue(ctx context.Context, db *gorm.DB, userid string, clue string) (result []model.File, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Where("user_uuid = ? AND filename LIKE ?", userid, "%"+clue+"%").Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (f *fileDAO) GetFileBucket(ctx context.Context, db *gorm.DB, userid, runid, taskid, filename string) (result string, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Select("bucket").Where("username = ? AND run_uuid = ? AND task_uuid = ? AND filename = ?", userid, runid, taskid, filename).Take(&result).Error; err != nil {
		return "", err
	}
	return result, nil
}

func (f *fileDAO) GetExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.File, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Where("expired_at > '0001-01-01 00:00:00' AND expired_at <= ?", date).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (f *fileDAO) GetExpiredFilesByRunID(ctx context.Context, db *gorm.DB, userid, runid string, date time.Time) (result []model.File, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND expired_at > '0001-01-01 00:00:00' AND expired_at <= ?", userid, runid, date).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil

}

func (f *fileDAO) GetFilesByOffsetWithFilters(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.File, err error) {
	thisDB := db.WithContext(ctx).Table("files")
	thisDB = thisDB.Where("username= ?", userid)

	if filter != "" {
		filterParamSlice := strings.Split(filter, ";")
		for _, v := range filterParamSlice {
			var itemKeyValue []string
			if strings.Contains(v, ">=") {
				itemKeyValue = strings.Split(v, ">=")
				thisDB = thisDB.Where(fmt.Sprintf("%s >= ?", itemKeyValue[0]), itemKeyValue[1])
			} else if strings.Contains(v, "<=") {
				itemKeyValue = strings.Split(v, "<=")
				thisDB = thisDB.Where(fmt.Sprintf("%s <= ?", itemKeyValue[0]), itemKeyValue[1])
			} else {
				itemKeyValue = strings.Split(v, "=")
				thisDB = thisDB.Where(fmt.Sprintf("%s LIKE ?", itemKeyValue[0]), "%"+itemKeyValue[1]+"%")
			}
		}
	}

	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Offset(offset).Limit(pageSize)

	if orderParam != "" {
		thisDB.Order(orderParam)
	} else {
		thisDB.Order("created_at desc")
	}
	err = thisDB.Find(&result).Error
	return result, err
}

func (f *fileDAO) GetFilesByOffsetWithFiltersForUserUploadFile(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.File, total int64, err error) {
	thisDB := db.WithContext(ctx).Table("files")
	thisDB = thisDB.Where("username= ? AND user_upload_name <> ?", userid, "")

	if filter != "" {
		filterParamSlice := strings.Split(filter, ";")
		for _, v := range filterParamSlice {
			var itemKeyValue []string
			if strings.Contains(v, ">=") {
				itemKeyValue = strings.Split(v, ">=")
				thisDB = thisDB.Where(fmt.Sprintf("%s >= ?", itemKeyValue[0]), itemKeyValue[1])
			} else if strings.Contains(v, "<=") {
				itemKeyValue = strings.Split(v, "<=")
				thisDB = thisDB.Where(fmt.Sprintf("%s <= ?", itemKeyValue[0]), itemKeyValue[1])
			} else {
				itemKeyValue = strings.Split(v, "=")
				thisDB = thisDB.Where(fmt.Sprintf("%s LIKE ?", itemKeyValue[0]), "%"+itemKeyValue[1]+"%")
			}
		}
	}

	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Count(&total).Offset(offset).Limit(pageSize)

	if orderParam != "" {
		thisDB.Order(orderParam).Order("sample_name")
	} else {
		thisDB.Order("Created_At desc, sample_name")
	}
	err = thisDB.Find(&result).Error
	return result, total, err
}

func (f *fileDAO) GetUserUploadFiles(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string, workflow_uuid string) (result []model.UserUploadFile, err error) {
	var query = `WITH upload(sample_name, total) AS (SELECT sample_name, COUNT(*) OVER() FROM files WHERE workflow_uuid = ? and username = ? GROUP BY sample_name LIMIT ? OFFSET ?)
																SELECT f.sample_name, f.username, f.local_path, f.user_upload_name, f.filename, f.filesize, f.bucket, f.created_at, upload.total  FROM files AS f 
																INNER JOIN upload ON f.sample_name = upload.sample_name
																WHERE username = ?
																ORDER BY created_at DESC
																`
	err = db.WithContext(ctx).Raw(query, workflow_uuid, userid, pageSize, pageSize*(pageToken-1), userid).Scan(&result).Error
	return
}

func (f *fileDAO) GetFilesByOffsetWithFiltersForSystemGenFile(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.File, err error) {
	thisDB := db.WithContext(ctx).Table("files")
	thisDB = thisDB.Where("username= ? AND user_upload_name= ?", userid, "")

	if filter != "" {
		filterParamSlice := strings.Split(filter, ";")
		for _, v := range filterParamSlice {
			var itemKeyValue []string
			if strings.Contains(v, ">=") {
				itemKeyValue = strings.Split(v, ">=")
				thisDB = thisDB.Where(fmt.Sprintf("%s >= ?", itemKeyValue[0]), itemKeyValue[1])
			} else if strings.Contains(v, "<=") {
				itemKeyValue = strings.Split(v, "<=")
				thisDB = thisDB.Where(fmt.Sprintf("%s <= ?", itemKeyValue[0]), itemKeyValue[1])
			} else {
				itemKeyValue = strings.Split(v, "=")
				thisDB = thisDB.Where(fmt.Sprintf("%s LIKE ?", itemKeyValue[0]), "%"+itemKeyValue[1]+"%")
			}
		}
	}

	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Offset(offset).Limit(pageSize)

	if orderParam != "" {
		thisDB.Order(orderParam)
	}
	err = thisDB.Find(&result).Error
	return result, err
}

func (f *fileDAO) GetFileNotUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.File, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Where("upload_success = ? AND username = ? AND run_uuid ?", false, userid, runid).Find(&result).Error; err != nil {
		return nil, err
	}

	return
}

func (f *fileDAO) GetFileUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.File, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Where("upload_success = ? AND username = ? AND run_uuid ?", true, userid, runid).Find(&result).Error; err != nil {
		return nil, err
	}

	return
}

func (f *fileDAO) GetFileByFilename(ctx context.Context, db *gorm.DB, userid string, runid string, taskid string, filename string) (result *model.File, err error) {
	result = new(model.File)
	//TODO(tuandn8) need to filter with userid
	if err := db.Model(&model.File{}).WithContext(ctx).Where("upload_success = ? AND run_uuid = ? AND task_uuid = ? AND filename = ?", false, runid, taskid, filename).Find(result).Error; err != nil {
		return nil, err
	}

	return
}

func (f *fileDAO) GetExpiredUploadFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.File, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Where("upload_expired_at > '0001-01-01 00:00:00' AND upload_expired_at <= ? AND upload_success = ? AND done_run = ?", date, false, true).Find(&result).Error; err != nil {
		return nil, err
	}

	return
}
func (f *fileDAO) UpdateFileAsUploaded(ctx context.Context, db *gorm.DB, userid, runid, taskid, bucket, filename string) (err error) {
	var file = &model.File{
		Bucket: bucket,
	}
	if err := db.Model(&model.File{}).WithContext(ctx).Where("username = ? AND run_uuid = ?, task_uuid = ?, filename = ?", userid, runid, taskid, filename).Updates(file).Error; err != nil {
		return err
	}

	return
}

func (f *fileDAO) SampleExisted(ctx context.Context, db *gorm.DB, userid, sample_name, workflow string) (existed bool, err error) {
	var file = &model.File{}
	if err := db.Model(&model.File{}).WithContext(ctx).Where("username = ?", userid).Where("sample_name = ?", sample_name).Where("workflow_uuid = ?", workflow).First(file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f *fileDAO) DeleteFileFromHardDisk(ctx context.Context, db *gorm.DB, userid, fileId string, removeFile bool) (err error) {
	var isGeneratedFile bool = true
	var generatedFile model.GeneratedFile
	err = db.WithContext(ctx).Where(&model.GeneratedFile{FileUUID: fileId, UserID: userid}).First(&generatedFile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isGeneratedFile = false
		} else {
			return err
		}
	}

	var isUploadedFile bool = true
	var uploadedFile model.UploadedFile
	err = db.WithContext(ctx).Where(&model.UploadedFile{FileUUID: fileId, UserID: userid}).First(&uploadedFile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isUploadedFile = false
		} else {
			return err
		}
	}

	if isGeneratedFile {
		err = utils.DeleteFile(generatedFile.Path)
		if err = db.WithContext(ctx).Where("file_uuid = ?", fileId).Delete(&model.GeneratedFile{}).Error; err != nil {
			return err
		}
	}

	if isUploadedFile {
		err = utils.DeleteFile(uploadedFile.Path)
		if err = db.WithContext(ctx).Where("file_uuid = ?", fileId).Delete(&model.UploadedFile{}).Error; err != nil {
			return err
		}
	}
	return
}
