package workflowdag

import (
	"fmt"
)

// AddOutputToInput add output regex & secondary file to step input
func AddOutputToInput(wf *WorkflowDAG) (err error) {
	var (
		stepMap = make(map[string]*Step)
	)

	for _, step := range wf.Steps {
		if _, ok := stepMap[step.WorkflowName]; ok {
			return fmt.Errorf("Duplicate step name: %s", step.ID)
		}
		stepMap[step.WorkflowName] = step
	}

	for _, step := range wf.Steps {
		for _, argument := range step.Arguments {
			if argument.Input == nil {
				continue
			}
			if argument.Input.From == "" {
				continue
			}
			if step, ok := stepMap[argument.Input.From]; ok {
				for _, stepOutput := range step.StepOutput {
					if stepOutput.Name == argument.Input.WorkflowName {
						argument.Input.Value = append(argument.Input.Value, stepOutput.Regex...)
						argument.Input.SecondaryFiles = append(argument.Input.SecondaryFiles, stepOutput.SecondaryFiles...)
					}
				}
				continue
			}

			return fmt.Errorf("Can not find step: %s. In step map: %v", argument.Input.From, stepMap)
		}
	}

	return nil
}
