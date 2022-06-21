package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"workflow/valkyrie/core"
	"workflow/valkyrie/model"
	"gorm.io/gorm"
)

type genFileDAO struct{}

type IgenFileDAO interface {
	SaveFile(ctx context.Context, db *gorm.DB, file *model.GeneratedFile) (err error)
	DeleteFile(ctx context.Context, db *gorm.DB, userid, runid, taskid, filename string) (err error)
	UpdateCloudExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error)
	UpdateHardDiskExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error)
	UpdateDoneRun(ctx context.Context, db *gorm.DB, userid, runid string) (err error)

	GetFileBucket(ctx context.Context, db *gorm.DB, fileUUID string) (result string, err error)

	// for download from disk
	GetFileByFilename(ctx context.Context, db *gorm.DB, userid string, runid string, taskid string, filename string) (result *model.GeneratedFile, err error)

	GetFileNotUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.GeneratedFile, err error)
	GetFileUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.GeneratedFile, err error)

	// For file deletion
	GetExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.GeneratedFile, err error)
	GetExpiredFilesByRunID(ctx context.Context, db *gorm.DB, userid, runid string, date time.Time) (result []model.GeneratedFile, err error)
	GetUploadExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.GeneratedFile, err error)

	UpdateFileAsUploaded(ctx context.Context, db *gorm.DB, bucket, fileUUID, path string) (err error)
	GetUserGeneratedFiles(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.GeneratedFile, total int64, err error)
}

func (g *genFileDAO) SaveFile(ctx context.Context, db *gorm.DB, file *model.GeneratedFile) (err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Create(file).Error; err != nil {
		return err
	}

	return nil
}

func (g *genFileDAO) DeleteFile(ctx context.Context, db *gorm.DB, userid, runid, taskid, filename string) (err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND task_uuid = ? AND filename = ?", userid, runid, taskid, filename).Delete(&model.GeneratedFile{}).Error; err != nil {
		return err
	}

	return nil
}

func (g *genFileDAO) UpdateCloudExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error) {
	mainConf := core.GetMainConfig()
	file := &model.GeneratedFile{
		ExpiredAt: expiredAt,
	}

	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND upload_success = ?", userid, runid, !mainConf.HardDiskOnly).Updates(file).Error; err != nil {
		return err
	}

	return nil
}

func (g *genFileDAO) UpdateHardDiskExpiredTime(ctx context.Context, db *gorm.DB, userid, runid string, expiredAt time.Time) (err error) {
	file := &model.GeneratedFile{
		UploadExpiredAt: expiredAt,
	}

	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND upload_success = ?", userid, runid, false).Updates(file).Error; err != nil {
		return err
	}

	return nil
}

func (g *genFileDAO) UpdateDoneRun(ctx context.Context, db *gorm.DB, userid, runid string) (err error) {
	file := &model.GeneratedFile{
		DoneRun: true,
	}

	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("username = ? AND run_uuid = ?", userid, runid).Updates(file).Error; err != nil {
		return err
	}

	return
}

func (g *genFileDAO) GetFileBucket(ctx context.Context, db *gorm.DB, fileUUID string) (result string, err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Select("bucket").Where("file_uuid=?", fileUUID).Take(&result).Error; err != nil {
		return "", err
	}
	return result, nil
}

func (g *genFileDAO) GetExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.GeneratedFile, err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("expired_at > '0001-01-01 00:00:00' AND expired_at <= ?", date).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (g *genFileDAO) GetExpiredFilesByRunID(ctx context.Context, db *gorm.DB, userid, runid string, date time.Time) (result []model.GeneratedFile, err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("username = ? AND run_uuid = ? AND expired_at > '0001-01-01 00:00:00' AND expired_at <= ?", userid, runid, date).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (g *genFileDAO) UpdateFileAsUploaded(ctx context.Context, db *gorm.DB, bucket, fileUUID, path string) (err error) {
	var file = &model.GeneratedFile{
		Bucket: bucket,
		Path:   path,
	}
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("file_uuid = ?", fileUUID).Updates(file).Error; err != nil {
		return err
	}

	return
}

func (g *genFileDAO) GetFileNotUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.GeneratedFile, err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("upload_success = ? AND username = ? AND run_uuid ?", false, userid, runid).Find(&result).Error; err != nil {
		return nil, err
	}

	return
}

func (g *genFileDAO) GetFileUploadedByRunID(ctx context.Context, db *gorm.DB, userid string, runid string) (result []model.GeneratedFile, err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("upload_success = ? AND username = ? AND run_uuid ?", true, userid, runid).Find(&result).Error; err != nil {
		return nil, err
	}

	return
}

func (g *genFileDAO) GetFileByFilename(ctx context.Context, db *gorm.DB, userid string, runid string, taskid string, filename string) (result *model.GeneratedFile, err error) {
	result = new(model.GeneratedFile)
	//TODO(tuandn8) need to filter with userid
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("upload_success = ? AND run_uuid = ? AND task_uuid = ? AND filename = ?", false, runid, taskid, filename).Find(result).Error; err != nil {
		return nil, err
	}

	return
}

func (g *genFileDAO) GetUploadExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.GeneratedFile, err error) {
	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("upload_expired_at > '0001-01-01 00:00:00' AND upload_expired_at <= ? AND upload_success = ? AND done_run = ?", date, false, true).Find(&result).Error; err != nil {
		return nil, err
	}

	return
}

func (u *genFileDAO) GetUserGeneratedFiles(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.GeneratedFile, total int64, err error) {

	thisDB := db.WithContext(ctx).Table("generated_files")
	thisDB = thisDB.Where("username= ?", userid)
	countQuery := db.WithContext(ctx).Model(&model.GeneratedFile{})
	countQuery = countQuery.Where("username= ?", userid)

	if filter != "" {
		filterParamSlice := strings.Split(filter, ";")
		for _, v := range filterParamSlice {
			var itemKeyValue []string
			if strings.Contains(v, ">=") {
				itemKeyValue = strings.Split(v, ">=")
				thisDB = thisDB.Where(fmt.Sprintf("%s >= ?", itemKeyValue[0]), itemKeyValue[1])
				countQuery = countQuery.Where(fmt.Sprintf("%s >= ?", itemKeyValue[0]), itemKeyValue[1])
			} else if strings.Contains(v, "<=") {
				itemKeyValue = strings.Split(v, "<=")
				thisDB = thisDB.Where(fmt.Sprintf("%s <= ?", itemKeyValue[0]), itemKeyValue[1])
				countQuery = countQuery.Where(fmt.Sprintf("%s <= ?", itemKeyValue[0]), itemKeyValue[1])
			} else {
				itemKeyValue = strings.Split(v, "=")
				thisDB = thisDB.Where(fmt.Sprintf("%s LIKE ?", itemKeyValue[0]), "%"+itemKeyValue[1]+"%")
				countQuery = countQuery.Where(fmt.Sprintf("%s LIKE ?", itemKeyValue[0]), "%"+itemKeyValue[1]+"%")
			}
		}
	}

	offset := (pageToken - 1) * pageSize
	thisDB = thisDB.Offset(offset).Limit(pageSize)

	if orderParam != "" {
		thisDB.Order(orderParam)
	}
	err = thisDB.Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	err = countQuery.Count(&total).Error

	return result, total, err
}

/*
func (g *genFileDAO) UpdateFileSampleInfo(ctx context.Context, db *gorm.DB, fileUUID, sampleUUID string, sampleIndex int) (err error) {
	var file = &model.GeneratedFile{
		SampleUUID:  sampleUUID,
		SampleIndex: sampleIndex,
	}

	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("file_uuid = ?", fileUUID).Updates(file).Error; err != nil {
		return err
	}

	return
}
*/
var (
	iGenFileDAO IgenFileDAO = &genFileDAO{}
)

func SetGeneratedFileDAO(v IgenFileDAO) {
	iGenFileDAO = v
	return
}

func GetGeneratedFileDAO() (dao IgenFileDAO) {
	return iGenFileDAO
}
