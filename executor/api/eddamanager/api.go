package eddamanager

import (
	"encoding/json"
	"strings"

	"github.com/vfluxus/dvergr/constant"
	"github.com/vfluxus/dvergr/httpdto/eddamanagerdto"
	"github.com/vfluxus/workflow-utils/model"
	"github.com/vfluxus/workflow/executor/api"
	"github.com/vfluxus/workflow/executor/core"
)

const (
	tasklogNew   = "/tasklog/"
	tasklogState = "/tasklog/state"
)

func NewTaskLog(taskUUID, podName, nameSpace, containerName, containerID, nameNode string) (ok bool, err error) {
	var data []byte
	var response = &model.Response{}
	containerIDraw := strings.ReplaceAll(containerID, "docker://", "")
	payload := &eddamanagerdto.NewTaskLogReq{
		TaskUUID:      taskUUID,
		MetadataName:  podName,
		NameSpace:     nameSpace,
		ContainerName: containerName,
		ContainerID:   containerIDraw,
		State:         constant.TaskStateRunning,
		NodeName:      nameNode,
	}
	url := core.GetMainConfig().EddaConfig.Addr

	if data, err = json.Marshal(payload); err != nil {
		return false, err
	}

	err = api.PostJSON(url+tasklogNew, false, "", data, response)
	if err != nil {
		return false, err
	}
	_, err = api.ProcessResponse(response)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateTaskLogState(taskUUID string) (ok bool, err error) {
	var data []byte
	var response = &model.Response{}
	payload := &eddamanagerdto.UpdateTaskLogStateReq{
		TaskUUID: taskUUID,
		State:    constant.TaskStateComplete,
	}
	url := core.GetMainConfig().EddaConfig.Addr

	if data, err = json.Marshal(payload); err != nil {
		return false, err
	}

	err = api.PostJSON(url+tasklogState, false, "", data, response)
	if err != nil {
		return false, err
	}
	_, err = api.ProcessResponse(response)
	if err != nil {
		return false, err
	}

	return true, nil
}
