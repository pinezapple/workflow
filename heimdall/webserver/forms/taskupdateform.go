package forms

import (
	"github.com/google/uuid"
	"github.com/vfluxus/heimdall/services/dto"
)

type UpdateTaskStatusForm struct {
	TaskIDs   []string    `json:"task_ids"`
	TaskUUIDs []uuid.UUID `json:"task_uuids"`
	Status    uint32      `json:"status"`
}

type UpdateTaskOutputLocationForm struct {
	TaskID         string   `json:"task_id"`
	OutputLocation []string `json:"output_location"`
}

type UpdateParamWithRegexFileSizeForm struct {
	TaskID         string                `json:"task_id"`
	FileSize       []int                 `json:"file_size"`
	ParamWithRegex []*dto.ParamWithRegex `json:"param_with_regex"`
}

type TaskUpdateResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Data interface{} `json:"data"`
}

func GetTaskUpdateResp(code int, message string, data interface{}) *TaskUpdateResponse {
	var (
		resp = &TaskUpdateResponse{
			Data: data,
		}
	)
	resp.Error.Code = code
	resp.Error.Message = message

	return resp
}
