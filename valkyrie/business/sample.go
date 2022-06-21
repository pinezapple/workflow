package business

import (
	"context"

	"github.com/google/uuid"
	"workflow/valkyrie/core"
	"workflow/valkyrie/dao"
	"workflow/valkyrie/model"
)

func CreateSample(ctx context.Context, req *model.CreateNewSampleReq) (err error) {
	db := core.GetDBObj()
	userID := ctx.Value("UserID").(string)
	sampleDAO := dao.GetSampleDAO()

	newSample := &model.Sample{
		SampleUUID:   uuid.New().String(),
		SampleName:   req.SampleName,
		DatasetUUID:  req.DatasetUUID,
		WorkflowUUID: req.WorkflowUUID,
		UserID:       userID,
	}
	_, err = sampleDAO.SaveSample(ctx, db, newSample, req.FileUUIDs)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return
}

func GetSampleByDataset(ctx context.Context, datasetUUID string, pageSize, pageToken int, filter, orderParam string) (sampleWithDetails []*model.SampleDetailResp, err error) {
	db := core.GetDBObj()
	sampleDAO := dao.GetSampleDAO()

	sample, err := sampleDAO.GetSampleByDataset(ctx, db, datasetUUID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	for i := 0; i < len(sample); i++ {
		var sampleWithDetail = &model.SampleDetailResp{
			SampleUUID:   sample[i].SampleUUID,
			DatasetUUID:  sample[i].DatasetUUID,
			UserID:       sample[i].UserID,
			SampleName:   sample[i].SampleName,
			WorkflowUUID: sample[i].WorkflowUUID,
			CreatedAt:    sample[i].CreatedAt,
		}
		detail, err := sampleDAO.GetFilesBySampleID(ctx, db, sample[i].SampleUUID)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		sampleWithDetail.SampleFiles = detail
		sampleWithDetails = append(sampleWithDetails, sampleWithDetail)
	}

	return sampleWithDetails, nil
}

func GetSampleByWorkflow(ctx context.Context, userID, workflowUUID string, pageSize, pageToken int, filter, orderParam string) (sampleWithDetails []*model.SampleDetailResp, total int64, err error) {
	db := core.GetDBObj()
	sampleDAO := dao.GetSampleDAO()

	sample, total, err := sampleDAO.GetSamplesByWorkflow(ctx, db, userID, workflowUUID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		logger.Error(err.Error())
		return nil, 0, err
	}

	for i := 0; i < len(sample); i++ {
		var sampleWithDetail = &model.SampleDetailResp{
			SampleUUID:   sample[i].SampleUUID,
			DatasetUUID:  sample[i].DatasetUUID,
			UserID:       sample[i].UserID,
			SampleName:   sample[i].SampleName,
			WorkflowUUID: sample[i].WorkflowUUID,
			CreatedAt:    sample[i].CreatedAt,
		}
		detail, err := sampleDAO.GetFilesBySampleID(ctx, db, sample[i].SampleUUID)
		if err != nil {
			logger.Error(err.Error())
			return nil, 0, err
		}
		sampleWithDetail.SampleFiles = detail
		sampleWithDetails = append(sampleWithDetails, sampleWithDetail)
	}

	return sampleWithDetails, total, nil
}

func GetSampleFilter(ctx context.Context, userID string, pageSize, pageToken int, filter, orderParam string) (sampleWithDetails []*model.SampleDetailResp, err error) {
	db := core.GetDBObj()
	sampleDAO := dao.GetSampleDAO()

	sample, err := sampleDAO.GetSamples(ctx, db, userID, pageSize, pageToken, filter, orderParam)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	for i := 0; i < len(sample); i++ {
		var sampleWithDetail = &model.SampleDetailResp{
			SampleUUID:   sample[i].SampleUUID,
			DatasetUUID:  sample[i].DatasetUUID,
			UserID:       sample[i].UserID,
			SampleName:   sample[i].SampleName,
			WorkflowUUID: sample[i].WorkflowUUID,
			CreatedAt:    sample[i].CreatedAt,
		}
		detail, err := sampleDAO.GetFilesBySampleID(ctx, db, sample[i].SampleUUID)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		sampleWithDetail.SampleFiles = detail
		sampleWithDetails = append(sampleWithDetails, sampleWithDetail)
	}

	return sampleWithDetails, nil
}

func GetSampleDetail(ctx context.Context, sampleUUID string) (sampleDetail *model.SampleDetailResp, err error) {
	db := core.GetDBObj()
	sampleDAO := dao.GetSampleDAO()

	sampleMeta, err := sampleDAO.GetSampleMetadataBySampleUUID(ctx, db, sampleUUID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	files, err := sampleDAO.GetFilesBySampleID(ctx, db, sampleUUID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	resp := &model.SampleDetailResp{
		SampleUUID:   sampleMeta.SampleUUID,
		DatasetUUID:  sampleMeta.DatasetUUID,
		UserID:       sampleMeta.UserID,
		SampleName:   sampleMeta.SampleName,
		WorkflowUUID: sampleMeta.WorkflowUUID,
		SampleFiles:  files,
		CreatedAt:    sampleMeta.CreatedAt,
	}

	return resp, nil
}

func CheckExistSample(ctx context.Context, sampleName, workflow, userid string) (ok bool, err error) {
	db := core.GetDBObj()
	sampleDAO := dao.GetSampleDAO()

	existed, err := sampleDAO.GetSampleMetadataBySampleName(ctx, db, sampleName, userid, workflow)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	} else {
		if existed.SampleUUID != "" {
			return true, nil
		} else {
			return false, nil
		}
	}
}
