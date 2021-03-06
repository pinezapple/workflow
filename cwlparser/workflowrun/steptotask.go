package workflowrun

import (
	"strings"

	"workflow/cwlparser/workflowdag"
)

// convertFromStepDAGToTask 1 step -> 1 task, 1 arguments -> 1 param with regex
func convertFromStepDAGToTask(step *workflowdag.Step, taskID string, taskName string) (task *Task, err error) {
	task = &Task{
		TaskID:          taskID,
		TaskName:        taskName,
		IsBoundary:      false,
		StepID:          step.ID,
		ScatterMethod:   step.ScatterMethod,
		ChildrenTasksID: step.ChildrenName, // for later replacement
		ParentTasksID:   step.ParentName,   // for later replacement
		Command:         strings.Join(step.BaseCommand, " ") + " ",
		DockerImage:     []string{step.DockerImage},
		OutputLocation:  nil,
		QueueLevel:      0,
		Status:          0,
	}
	var (
		taskParamWithRegex = make([]*ParamWithRegex, len(step.Arguments))
	)

	for argIndex := range step.Arguments {
		if step.Arguments[argIndex].Input == nil {
			newParamWithRegex := &ParamWithRegex{
				Prefix: step.Arguments[argIndex].Prefix,
			}
			taskParamWithRegex[argIndex] = newParamWithRegex
			continue
		}

		var (
			from []string
		)
		if step.Arguments[argIndex].Input.From != "" {
			from = append(from, step.Arguments[argIndex].Input.From)
		}
		newParamWithRegex := &ParamWithRegex{
			From:           from,
			Scatter:        step.Arguments[argIndex].Input.Scatter,
			SecondaryFiles: step.Arguments[argIndex].Input.SecondaryFiles,
			Regex:          step.Arguments[argIndex].Input.Value,
			Prefix:         step.Arguments[argIndex].Prefix,
		}

		// for i := range newParamWithRegex.Regex {
		// 	newParamWithRegex.Files = append(newParamWithRegex.Files, &FilteredFiles{
		// 		Filepath: newParamWithRegex.Regex[i],
		// 	})
		// }

		taskParamWithRegex[argIndex] = newParamWithRegex
	}

	task.ParamsWithRegex = taskParamWithRegex

	// build regex output
	for stepOutputIndex := range step.StepOutput {
		for regexIndex := range step.StepOutput[stepOutputIndex].Regex {
			task.OutputRegex = append(task.OutputRegex, step.StepOutput[stepOutputIndex].Regex[regexIndex])
		}
		for secondFilesIndex := range step.StepOutput[stepOutputIndex].SecondaryFiles {
			task.Output2ndFiles = append(task.Output2ndFiles, step.StepOutput[stepOutputIndex].SecondaryFiles[secondFilesIndex])
		}
	}

	return task, nil
}
