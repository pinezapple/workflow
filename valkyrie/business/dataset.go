package business

import (
	"context"

	"github.com/google/uuid"
	"workflow/valkyrie/core"
	"workflow/valkyrie/dao"
	"workflow/valkyrie/model"
)

func CreateDataset(ctx context.Context, req *model.CreateNewDatasetReq) (data interface{}, err error) {
	db := core.GetDBObj()
	userID := "tungnt99"
	datasetDAO := dao.GetDatasetDAO()

	newDataset := &model.Dataset{
		DatasetUUID: uuid.New().String(),
		DatasetName: req.DatasetName,
		UserID:      userID,
	}

	_, err = datasetDAO.SaveDataset(ctx, db, newDataset)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func GetDatasetWithFilter(ctx context.Context, pageSize, pageToken int, filter, orderParam string) (data []model.Dataset, err error) {
	db := core.GetDBObj()
	datasetDAO := dao.GetDatasetDAO()
	userID := "tungnt99"

	dataset, err := datasetDAO.GetDatasetByOffset(ctx, db, userID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return dataset, nil
}
