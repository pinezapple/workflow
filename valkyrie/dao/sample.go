package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/vfluxus/valkyrie/model"
	"gorm.io/gorm"
)

type sampleDAO struct{}

type IsampleDAO interface {
	SaveSample(ctx context.Context, db *gorm.DB, sample *model.Sample, fileUUID []string) (sampleID string, err error)
	AddFileToExistedSample(ctx context.Context, db *gorm.DB, sample *model.Sample, fileUUID string, index int) (err error)

	DeleteSample(ctx context.Context, db *gorm.DB, sampleUUID, userID string) (err error)

	GetSampleMetadataBySampleName(ctx context.Context, db *gorm.DB, sampleName, userID, workflowID string) (sample *model.Sample, err error)
	GetSampleMetadataBySampleUUID(ctx context.Context, db *gorm.DB, sampleUUID string) (sample *model.Sample, err error)

	GetAllSampleByDataset(ctx context.Context, db *gorm.DB, datasetUUID string) (sample []model.Sample, err error)

	GetSampleByDataset(ctx context.Context, db *gorm.DB, datasetUUID string, pageSize int, pageToken int, filter string, orderParam string) (sample []model.Sample, err error)

	GetSamples(ctx context.Context, db *gorm.DB, userID string, pageSize int, pageToken int, filter string, orderParam string) (sample []model.Sample, err error)

	GetSampleFileIndex(ctx context.Context, db *gorm.DB, sampleUUID string) (content []model.SampleContent, err error)

	GetSamplesByWorkflow(ctx context.Context, db *gorm.DB, userID, workflowUUID string, pageSize int, pageToken int, filter string, orderParam string) (sample []model.Sample, total int64, err error)

	GetFilesBySampleID(ctx context.Context, db *gorm.DB, sampleUUID string) (files []model.MutualFile, err error)
}

func (s *sampleDAO) SaveSample(ctx context.Context, db *gorm.DB, sample *model.Sample, fileUUID []string) (sampleID string, err error) {

	tx := db.WithContext(ctx).Begin()

	if err := tx.Model(&model.Sample{}).WithContext(ctx).Create(sample).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	for i := 0; i < len(fileUUID); i++ {
		var sampleContent = &model.SampleContent{
			FileUUID:    fileUUID[i],
			SampleUUID:  sample.SampleUUID,
			SampleIndex: i,
		}
		if err := tx.Model(&model.SampleContent{}).WithContext(ctx).Create(sampleContent).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()

	return sample.SampleUUID, nil
}

func (s *sampleDAO) AddFileToExistedSample(ctx context.Context, db *gorm.DB, sample *model.Sample, fileUUID string, index int) (err error) {
	var sampleContent = &model.SampleContent{
		FileUUID:    fileUUID,
		SampleUUID:  sample.SampleUUID,
		SampleIndex: index,
	}
	if err := db.Model(&model.SampleContent{}).WithContext(ctx).Create(sampleContent).Error; err != nil {
		return err
	}
	return
}

//TODO: implement this
func (s *sampleDAO) DeleteSample(ctx context.Context, db *gorm.DB, sampleUUID, userID string) (err error) {

	tx := db.WithContext(ctx).Begin()

	if err := tx.Model(&model.Sample{}).WithContext(ctx).Where("sample_uuid= ? AND username = ?", sampleUUID, userID).Delete(&model.Sample{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.SampleContent{}).WithContext(ctx).Where("sample_uuid=?", sampleUUID).Delete(&model.SampleContent{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (s *sampleDAO) GetAllSampleByDataset(ctx context.Context, db *gorm.DB, datasetUUID string) (sample []model.Sample, err error) {
	if err := db.Model(&model.Sample{}).WithContext(ctx).Where("dataset_uuid = ?", datasetUUID).Find(&sample).Error; err != nil {
		return nil, err
	}

	return
}

func (s *sampleDAO) GetSampleByDataset(ctx context.Context, db *gorm.DB, datasetUUID string, pageSize int, pageToken int, filter string, orderParam string) (result []model.Sample, err error) {
	thisDB := db.WithContext(ctx).Table("samples")
	thisDB = thisDB.Where("dataset_uuid= ?", datasetUUID)

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

func (s *sampleDAO) GetSamplesByWorkflow(ctx context.Context, db *gorm.DB, userID, workflowUUID string, pageSize int, pageToken int, filter string, orderParam string) (result []model.Sample, total int64, err error) {
	thisDB := db.WithContext(ctx).Table("samples")
	thisDB = thisDB.Where("username= ?", userID)
	thisDB = thisDB.Where("workflow_uuid= ?", workflowUUID)

	countQuery := db.WithContext(ctx).Model(&model.Sample{})
	countQuery = countQuery.Where("username= ?", userID)
	countQuery = countQuery.Where("workflow_uuid= ?", workflowUUID)

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
	} else {
		thisDB.Order("created_at desc")
	}
	err = thisDB.Find(&result).Error
	if err != nil {
		return nil, 0, err
	}

	err = countQuery.Count(&total).Error
	return result, total, err
}

func (s *sampleDAO) GetSamples(ctx context.Context, db *gorm.DB, userID string, pageSize int, pageToken int, filter string, orderParam string) (result []model.Sample, err error) {
	thisDB := db.WithContext(ctx).Table("samples")
	thisDB = thisDB.Where("username= ?", userID)

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
func (s *sampleDAO) GetSampleMetadataBySampleName(ctx context.Context, db *gorm.DB, sampleName, userID, workflowID string) (sample *model.Sample, err error) {
	sample = new(model.Sample)
	if err := db.Model(&model.Sample{}).WithContext(ctx).Where("sample_name= ? AND username = ? AND workflow_uuid = ?", sampleName, userID, workflowID).Find(sample).Error; err != nil {
		return nil, err
	}

	return
}

func (s *sampleDAO) GetSampleMetadataBySampleUUID(ctx context.Context, db *gorm.DB, sampleUUID string) (sample *model.Sample, err error) {
	sample = new(model.Sample)
	if err := db.Model(&model.Sample{}).WithContext(ctx).Where("sample_uuid= ?", sampleUUID).Find(sample).Error; err != nil {
		return nil, err
	}

	return

}

func (s *sampleDAO) GetSampleFileIndex(ctx context.Context, db *gorm.DB, sampleUUID string) (content []model.SampleContent, err error) {
	if err := db.Model(&model.SampleContent{}).WithContext(ctx).Where("sample_uuid = ?", sampleUUID).Order("sample_index DESC").Find(&content).Error; err != nil {
		return nil, err
	}
	return
}

func (s *sampleDAO) GetFilesBySampleID(ctx context.Context, db *gorm.DB, sampleUUID string) (files []model.MutualFile, err error) {
	var uploadedFiles []model.MutualFile
	var generatedFiles []model.MutualFile

	var queryUploadedFile = "SELECT uf.file_uuid, uf.username, temp.sample_index, uf.bucket, uf.path, uf.filename, uf.filesize, uf.created_at FROM uploaded_files as uf JOIN (SELECT * FROM sample_contents as sc WHERE sc.sample_uuid = ?) as temp ON temp.file_uuid = uf.file_uuid ORDER BY temp.sample_index"

	var queryGeneratedFile = "SELECT uf.file_uuid, uf.username, temp.sample_index, uf.bucket, uf.path, uf.filename, uf.filesize, uf.created_at FROM generated_files as uf JOIN (SELECT * FROM sample_contents as sc WHERE sc.sample_uuid = ?) as temp ON temp.file_uuid = uf.file_uuid ORDER BY temp.sample_index"

	if err := db.WithContext(ctx).Raw(queryUploadedFile, sampleUUID).Scan(&uploadedFiles).Error; err != nil {
		return nil, err
	}

	if err := db.WithContext(ctx).Raw(queryGeneratedFile, sampleUUID).Scan(&generatedFiles).Error; err != nil {
		return nil, err
	}

	/*
		if err := db.Model(&model.UploadedFile{}).Select("file_uuid, sample_uuid, sample_index, user_id, bucket, path, file_name, file_size, created_at").WithContext(ctx).Where("sample_uuid= ?", sampleUUID).Order("sample_index").Find(&uploadedFiles).Error; err != nil {
			return nil, err
		}
		if err := db.Model(&model.GeneratedFile{}).Select("file_uuid, sample_uuid, sample_index, user_id, bucket, path, file_name, file_size, created_at").WithContext(ctx).Where("sample_uuid= ?", sampleUUID).Order("sample_index").Find(&generatedFiles).Error; err != nil {
			return nil, err
		}
	*/
	if len(generatedFiles) == 0 {
		return uploadedFiles, nil
	}
	if len(uploadedFiles) == 0 {
		return uploadedFiles, nil
	}
	p1 := 0
	p2 := 0
	for i := 0; i < len(generatedFiles)+len(uploadedFiles); i++ {
		if generatedFiles[p1].SampleIndex == i {
			files = append(files, generatedFiles[p1])
			p1++
		} else if uploadedFiles[p2].SampleIndex == i {
			files = append(files, uploadedFiles[p2])
			p2++
		}
	}
	return
}

var (
	iSampleDAO IsampleDAO = &sampleDAO{}
)

func SetSampleDAO(v IsampleDAO) {
	iSampleDAO = v
	return
}

func GetSampleDAO() (dao IsampleDAO) {
	return iSampleDAO
}
