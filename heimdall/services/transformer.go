package services

import (
	"context"
	"encoding/json"
	"fmt"

	"workflow/cwlparser"
	"workflow/cwlparser/workflowcwl"
	"workflow/heimdall/services/dto"
)

// Transformer interface defines transform methods
type Transformer interface {
	Transform(ctx context.Context, transformRequest dto.TransformRequest) (dto.TransformResponse, error)
}

// GetTransformerService returns transformer service implement
func GetTransformerService() Transformer {
	return transformerService{}
}

type transformerService struct{}

// TransformRun convert from workflow + params to run and tasks
func (c transformerService) Transform(ctx context.Context, transformRequest dto.TransformRequest) (dto.TransformResponse, error) {
	resp := new(dto.TransformResponse)
	cwlHTTPForm := &workflowcwl.HttpCWLForm{
		Name:    transformRequest.Name,
		Content: []byte(transformRequest.Content),
		Steps:   make([]*workflowcwl.HttpStepForm, 0, len(transformRequest.Steps)),
	}

	// add steps to cwlHTTPForm
	for i := range transformRequest.Steps {
		newHTTPCWLStep := &workflowcwl.HttpStepForm{
			Name:    transformRequest.Steps[i].Name,
			Content: []byte(transformRequest.Steps[i].Content),
		}
		cwlHTTPForm.Steps = append(cwlHTTPForm.Steps, newHTTPCWLStep)
	}

	// parse CWL to struct
	wfCWL, err := cwlparser.ParseCWLInMem(cwlHTTPForm)
	if err != nil {
		return *resp, fmt.Errorf("Transform - parse cwl in mem error: %v", err)
	}

	// create DAG
	wfDAG, err := cwlparser.CreateWorkflowDAG(wfCWL)
	if err != nil {
		return *resp, fmt.Errorf("Transform - create DAG error: %v", err)
	}

	// create run
	run, err := cwlparser.CreateRunFromWorkflow(wfDAG, transformRequest.Params, transformRequest.UserName, transformRequest.RunIndex)
	if err != nil {
		return *resp, fmt.Errorf("Transform - create RUN error: %v", err)
	}

	// create resp
	runDTO := dto.Run{
		RunID:    run.RunIndex,
		RunName:  run.RunName,
		UserName: run.UserName,
		Tasks:    make([]*dto.Task, 0, len(run.Tasks)),
	}
	for tIDX := range run.Tasks {
		dtoTask := &dto.Task{
			TaskID:          run.Tasks[tIDX].TaskID,
			TaskName:        run.Tasks[tIDX].TaskName,
			IsBoundary:      run.Tasks[tIDX].IsBoundary,
			StepName:        run.Tasks[tIDX].StepID,
			UserName:        run.Tasks[tIDX].UserName,
			Command:         run.Tasks[tIDX].Command,
			ParamsWithRegex: make([]*dto.ParamWithRegex, 0, len(run.Tasks[tIDX].ParamsWithRegex)),
			OutputRegex:     run.Tasks[tIDX].OutputRegex,
			Output2ndFiles:  run.Tasks[tIDX].Output2ndFiles,
			ParentTasksID:   run.Tasks[tIDX].ParentTasksID,
			ChildrenTasksID: run.Tasks[tIDX].ChildrenTasksID,
			OutputLocation:  run.Tasks[tIDX].OutputLocation,
			DockerImage:     run.Tasks[tIDX].DockerImage,
		}

		for pIDX := range run.Tasks[tIDX].ParamsWithRegex {
			paramRegex := &dto.ParamWithRegex{
				Scatter:        run.Tasks[tIDX].ParamsWithRegex[pIDX].Scatter,
				From:           run.Tasks[tIDX].ParamsWithRegex[pIDX].From,
				SecondaryFiles: run.Tasks[tIDX].ParamsWithRegex[pIDX].SecondaryFiles,
				Regex:          run.Tasks[tIDX].ParamsWithRegex[pIDX].Regex,
				Files:          make([]*dto.FilteredFiles, 0, len(run.Tasks[tIDX].ParamsWithRegex[pIDX].Files)),
				Prefix:         run.Tasks[tIDX].ParamsWithRegex[pIDX].Prefix,
			}

			for fIDX := range run.Tasks[tIDX].ParamsWithRegex[pIDX].Files {
				filterFiles := &dto.FilteredFiles{
					Filepath: run.Tasks[tIDX].ParamsWithRegex[pIDX].Files[fIDX].Filepath,
					Filesize: run.Tasks[tIDX].ParamsWithRegex[pIDX].Files[fIDX].Filesize,
				}
				paramRegex.Files = append(paramRegex.Files, filterFiles)
			}

			dtoTask.ParamsWithRegex = append(dtoTask.ParamsWithRegex, paramRegex)
		}

		runDTO.Tasks = append(runDTO.Tasks, dtoTask)
	}

	resp.Data = runDTO
	// marshal indent to response
	respB, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return *resp, fmt.Errorf("Transform - marshal indent error: %v", err)
	}
	logger.Debugf("Transformer response: \n%s \n", string(respB))

	return *resp, nil
}
