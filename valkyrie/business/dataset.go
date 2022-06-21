package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/vfluxus/valkyrie/core"
	"github.com/vfluxus/valkyrie/dao"
	"github.com/vfluxus/valkyrie/model"
)

func CreateDataset(ctx context.Context, req *model.CreateNewDatasetReq) (data interface{}, err error) {
	db := core.GetDBObj()
	userID := ctx.Value("UserID").(string)
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
	userID := ctx.Value("UserID").(string)

	dataset, err := datasetDAO.GetDatasetByOffset(ctx, db, userID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return dataset, nil
}
