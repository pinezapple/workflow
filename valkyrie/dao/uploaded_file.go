package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vfluxus/valkyrie/core"
	"github.com/vfluxus/valkyrie/model"
	"gorm.io/gorm"
)

type uploadedFileDAO struct{}

type IuploadedFileDAO interface {
	SaveFile(ctx context.Context, db *gorm.DB, file *model.UploadedFile) (fileUUID string, err error)
	DeleteFile(ctx context.Context, db *gorm.DB, fileUUID, userID string, save bool) (err error)

	UpdateCloudExpiredTime(ctx context.Context, db *gorm.DB, fileUUID string, expiredAt time.Time) (err error)
	UpdateHardDiskExpiredTime(ctx context.Context, db *gorm.DB, fileUUID string, expiredAt time.Time) (err error)

	GetFileBucket(ctx context.Context, db *gorm.DB, fileUUID string) (result string, err error)

	// TODO: maybe later, for download from disk
	//GetFileByFilename(ctx context.Context, db *gorm.DB, userid string, runid string, taskid string, filename string) (result *model.GeneratedFile, err error)

	GetUserUploadFiles(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.UploadedFile, total int64, err error)

	GetUserFiles(ctx context.Context, db *gorm.DB, userid string, filter string, orderParam string) (result []model.DataFile, err error)

	// For file deletion
	GetExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.UploadedFile, err error)
	GetUploadExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.UploadedFile, err error)

	UpdateFileAsUploaded(ctx context.Context, db *gorm.DB, bucket, fileUUID, path string) (err error)
	UpdateFileProjectPath(ctx context.Context, db *gorm.DB, pathFiles *model.UpdatePathFiles, userID string) (err error)
}

func (u *uploadedFileDAO) SaveFile(ctx context.Context, db *gorm.DB, file *model.UploadedFile) (fileUUID string, err error) {
	if err := db.Model(&model.UploadedFile{}).WithContext(ctx).Create(file).Error; err != nil {
		return "", err
	}

	return file.FileUUID, nil
}

func (u *uploadedFileDAO) DeleteFile(ctx context.Context, db *gorm.DB, fileUUID, userID string, save bool) (err error) {
	// update field
	if save {
		var (
			deleteFile = &model.UploadedFile{
				FileUUID: fileUUID,
				UserID:   userID,
				Deleted:  true,
			}
		)

		if err := db.Model(&model.UploadedFile{}).WithContext(ctx).Updates(deleteFile).Error; err != nil {
			return err
		}
	}

	// if not save, then run delete query
	if err := db.Model(&model.File{}).WithContext(ctx).Where("file_uuid = ? AND username = ?", fileUUID, userID).Delete(&model.Bucket{}).Error; err != nil {
		return err
	}

	return nil
}

func (u *uploadedFileDAO) UpdateCloudExpiredTime(ctx context.Context, db *gorm.DB, fileUUID string, expiredAt time.Time) (err error) {
	mainConf := core.GetMainConfig()
	file := &model.UploadedFile{
		ExpiredAt: expiredAt,
	}

	if err := db.Model(&model.UploadedFile{}).WithContext(ctx).Where("file_uuid = ? AND upload_success = ?", fileUUID, !mainConf.HardDiskOnly).Updates(file).Error; err != nil {
		return err
	}

	return nil
}

func (u *uploadedFileDAO) UpdateHardDiskExpiredTime(ctx context.Context, db *gorm.DB, fileUUID string, expiredAt time.Time) (err error) {
	file := &model.UploadedFile{
		UploadExpiredAt: expiredAt,
	}

	if err := db.Model(&model.UploadedFile{}).WithContext(ctx).Where("file_uuid= ? AND upload_success = ?", fileUUID, false).Updates(file).Error; err != nil {
		return err
	}

	return nil
}

func (u *uploadedFileDAO) GetFileBucket(ctx context.Context, db *gorm.DB, fileUUID string) (result string, err error) {
	if err := db.Model(&model.File{}).WithContext(ctx).Select("bucket").Where("file_uuid=?", fileUUID).Take(&result).Error; err != nil {
		return "", err
	}
	return result, nil
}

func (u *uploadedFileDAO) GetUserUploadFiles(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (result []model.UploadedFile, total int64, err error) {

	thisDB := db.WithContext(ctx).Table("uploaded_files")
	thisDB = thisDB.Where("username= ?", userid)
	countQuery := db.WithContext(ctx).Model(&model.UploadedFile{})
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

func (u *uploadedFileDAO) GetUserFiles(ctx context.Context, db *gorm.DB, userid string, filter string, orderParam string) (result []model.DataFile, err error) {
	var projectPath string = "/"
	var projectID string = ""

	if filter != "" {
		filterParamSlice := strings.Split(filter, ";")
		for _, v := range filterParamSlice {
			var itemKeyValue []string
			if strings.Contains(v, "project_path") {
				itemKeyValue = strings.Split(v, "=")
				projectPath = itemKeyValue[1]
			} else if strings.Contains(v, "project_id") {
				itemKeyValue = strings.Split(v, "project_id=")
				projectID = itemKeyValue[1]
			}
		}
	}

	// Query uploaded
	var uploadedFiles []model.UploadedFile
	// thisDB := db.WithContext(ctx).Where(&model.UploadedFile{UserID: userid})
	thisDB := db.WithContext(ctx).Where(&model.UploadedFile{ProjectPath: projectPath})
	if projectID != "" {
		thisDB.Where(&model.UploadedFile{ProjectID: projectID})
	}

	err = thisDB.Find(&uploadedFiles).Error
	if err != nil {
		return nil, err
	}

	// Query generated files
	var generatedFiles []model.GeneratedFile
	// thisDB = db.WithContext(ctx).Where(&model.GeneratedFile{UserID: userid})
	thisDB = db.WithContext(ctx).Where(&model.GeneratedFile{ProjectPath: projectPath})
	if projectID != "" {
		thisDB.Where(&model.GeneratedFile{ProjectID: projectID})
	}

	// TODO(tuandn8) Need to add filter by project id
	err = thisDB.Find(&generatedFiles).Error
	if err != nil {
		return nil, err
	}

	// Fuse all files
	for _, uploadedFile := range uploadedFiles {
		result = append(result, model.DataFile{
			ID:          uploadedFile.FileUUID,
			ProjectID:   uploadedFile.ProjectID,
			ProjectPath: uploadedFile.ProjectPath,
			Path:        uploadedFile.Path,
			Filename:    uploadedFile.Filename,
			Filesize:    uploadedFile.Filesize,
			Owner:       uploadedFile.UserID,
			CreatedAt:   uploadedFile.CreatedAt,
		})
	}

	for _, generatedFile := range generatedFiles {
		result = append(result, model.DataFile{
			ID:          generatedFile.FileUUID,
			ProjectID:   generatedFile.ProjectID,
			ProjectPath: generatedFile.ProjectPath,
			Path:        generatedFile.Path,
			Filename:    generatedFile.Filename,
			Filesize:    generatedFile.Filesize,
			Owner:       generatedFile.UserID,
			CreatedAt:   generatedFile.CreatedAt,
		})
	}

	return result, err
}

func (u *uploadedFileDAO) GetExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.UploadedFile, err error) {
	if err := db.Model(&model.UploadedFile{}).WithContext(ctx).Where("expired_at > '0001-01-01 00:00:00' AND expired_at <= ?", date).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (u *uploadedFileDAO) GetUploadExpiredFiles(ctx context.Context, db *gorm.DB, date time.Time) (result []model.UploadedFile, err error) {
	if err := db.Model(&model.UploadedFile{}).WithContext(ctx).Where("upload_expired_at > '0001-01-01 00:00:00' AND upload_expired_at <= ? AND upload_success = ?", date, false).Find(&result).Error; err != nil {
		return nil, err
	}

	return
}

func (u *uploadedFileDAO) UpdateFileAsUploaded(ctx context.Context, db *gorm.DB, bucket, fileUUID, path string) (err error) {
	var file = &model.UploadedFile{
		Bucket: bucket,
		Path:   path,
	}
	if err := db.Model(&model.UploadedFile{}).WithContext(ctx).Where("file_uuid = ?", fileUUID).Updates(file).Error; err != nil {
		return err
	}

	return
}

/*
func (u *uploadedFileDAO) UpdateFileSampleInfo(ctx context.Context, db *gorm.DB, fileUUID, sampleUUID string, sampleIndex int) (err error) {
	var file = &model.UploadedFile{
		SampleUUID:  sampleUUID,
		SampleIndex: sampleIndex,
	}

	if err := db.Model(&model.GeneratedFile{}).WithContext(ctx).Where("file_uuid = ?", fileUUID).Updates(file).Error; err != nil {
		return err
	}

	return
}
*/

func (u *uploadedFileDAO) UpdateFileProjectPath(ctx context.Context, db *gorm.DB, uploadPathFiles *model.UpdatePathFiles, userID string) (err error) {
	tx := db.WithContext(ctx).Begin()
	tx.SavePoint("sp2")

	for _, pathFile := range uploadPathFiles.PathFiles {

		file := &model.UploadedFile{
			ProjectPath: pathFile.ProjectPath,
		}

		generated_file := &model.GeneratedFile{
			ProjectPath: pathFile.ProjectPath,
		}

		if err = tx.Model(&model.UploadedFile{}).WithContext(ctx).Where(&model.UploadedFile{FileUUID: pathFile.FileID, UserID: userID}).Updates(file).Error; err != nil {
			tx.RollbackTo("sp2")
		}

		if err = tx.Model(&model.GeneratedFile{}).WithContext(ctx).Where(&model.GeneratedFile{FileUUID: pathFile.FileID, UserID: userID}).Updates(generated_file).Error; err != nil {
			tx.RollbackTo("sp2")
		}
	}

	return tx.Commit().Error
}

var (
	iUploadedFileDAO IuploadedFileDAO = &uploadedFileDAO{}
)

func SetUploadedFileDAO(v IuploadedFileDAO) {
	iUploadedFileDAO = v
	return
}

func GetUploadedFileDAO() (dao IuploadedFileDAO) {
	return iUploadedFileDAO
}
