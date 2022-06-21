package scheduler

import (
	"encoding/json"

	"github.com/vfluxus/workflow-utils/model"
	"github.com/vfluxus/workflow/executor/api"
	"github.com/vfluxus/workflow/executor/core"
	executorModel "github.com/vfluxus/workflow/executor/model"
)

const (
	selectTask                = "/selectTasks"
	executorResumeWithoutTask = "/executor/resume/new"
)

func SelectTask(cpu, ram int64) (ok bool, err error) {
	var data []byte
	var response = &model.Response{}
	mainConf := core.GetMainConfig()
	payload := &executorModel.SelectTaskReq{
		ExecutorName: mainConf.ServiceName,
		CPUpoint:     cpu,
		RAMpoint:     ram,
	}
	url := core.GetMainConfig().SchedulerConfig.Addr

	if data, err = json.Marshal(payload); err != nil {
		return false, err
	}

	err = api.PostJSON(url+selectTask, false, "", data, response)
	if err != nil {
		return false, err
	}
	_, err = api.ProcessResponse(response)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ExecutorResumeWithoutTask() (ok bool, err error) {
	var data []byte
	var response = &model.Response{}
	payload := &executorModel.ExecutorResumeWithoutTaskReq{
		ExecutorName: core.GetMainConfig().ServiceName,
	}
	url := core.GetMainConfig().SchedulerConfig.Addr

	if data, err = json.Marshal(payload); err != nil {
		return false, err
	}

	err = api.PostJSON(url+executorResumeWithoutTask, false, "", data, response)
	if err != nil {
		return false, err
	}
	_, err = api.ProcessResponse(response)
	if err != nil {
		return false, err
	}

	return true, nil
}
