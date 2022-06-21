package dao

import (
	"context"
	"fmt"
	"strings"

	"workflow/valkyrie/model"
	"gorm.io/gorm"
)

type datasetDAO struct{}

type IdatasetDAO interface {
	SaveDataset(ctx context.Context, db *gorm.DB, ds *model.Dataset) (dsID string, err error)
	DeleteDataset(ctx context.Context, db *gorm.DB, dsID, userid string) (err error)
	GetAllDataset(ctx context.Context, db *gorm.DB, userid string) (sample []model.Dataset, err error)
	GetDatasetByOffset(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (ds []model.Dataset, err error)
}

func (d *datasetDAO) SaveDataset(ctx context.Context, db *gorm.DB, ds *model.Dataset) (dsID string, err error) {
	if err := db.Model(&model.Dataset{}).WithContext(ctx).Create(ds).Error; err != nil {
		return "", err
	}

	return ds.DatasetUUID, nil
}

// TODO: implement this
func (d *datasetDAO) DeleteDataset(ctx context.Context, db *gorm.DB, dsID, userID string) (err error) {
	tx := db.WithContext(ctx).Begin()
	if err := tx.Model(&model.Dataset{}).WithContext(ctx).Where("sample_uuid= ? AND username= ?", dsID, userID).Delete(&model.Dataset{}).Error; err != nil {
		return err
	}
	//if err := tx.Model(&model.Sample).WithContext(ctx).Where("dataset_uuid = ?", dsID)

	return
}

func (d *datasetDAO) GetAllDataset(ctx context.Context, db *gorm.DB, userid string) (sample []model.Dataset, err error) {
	return
}

func (d *datasetDAO) GetDatasetByOffset(ctx context.Context, db *gorm.DB, userid string, pageSize int, pageToken int, filter string, orderParam string) (ds []model.Dataset, err error) {
	thisDB := db.WithContext(ctx).Table("datasets")
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
	err = thisDB.Find(&ds).Error
	return ds, err
}

var (
	iDatasetDAO IdatasetDAO = &datasetDAO{}
)

func SetDatasetDAO(v IdatasetDAO) {
	iDatasetDAO = v
	return
}

func GetDatasetDAO() (dao IdatasetDAO) {
	return iDatasetDAO
}
