package heimdall

import (
	"encoding/json"

	"github.com/vfluxus/valkyrie/api"
	"github.com/vfluxus/valkyrie/core"
	valkyrieModel "github.com/vfluxus/valkyrie/model"
	"github.com/vfluxus/workflow-utils/model"
)

const (
	uploadFileDone = ""
)

type UploadFileDoneReq struct {
	TaskUUID     string                          `json:"task_uuid"`
	SuccessFiles []*valkyrieModel.FileUploadResp `json:"success_files"`
	FailFiles    []*valkyrieModel.FileUploadResp `json:"fail_files"`
}

func UploadFileDoneStatus(taskUUID string, success, fail []*valkyrieModel.FileUploadResp) (err error) {
	var data []byte
	var response = &model.Response{}
	url := core.GetMainConfig().HeimdallConfig.Addr

	req := &UploadFileDoneReq{
		TaskUUID:     taskUUID,
		SuccessFiles: success,
		FailFiles:    fail,
	}
	if data, err = json.Marshal(req); err != nil {
		return
	}

	err = api.PostJSON(url+uploadFileDone, false, "", data, response)
	if err != nil {
		return err
	}

	_, err = api.ProcessResponse(response)
	if err != nil {
		return err
	} else {
		return nil
	}

}
