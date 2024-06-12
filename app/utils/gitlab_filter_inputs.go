package utils

import "github.com/mubashiroliyantakath/gitlab-metadata/app/interfaces"

// Filter a list of interface.Inputs based on GitLab environment variables
func filterInputs(inputs []interfaces.Input) []interfaces.Input {
	var filteredInputs []interfaces.Input

	for _, input := range inputs {
		if input.GitlabFilterPassed() {
			filteredInputs = append(filteredInputs, input)
		}
	}

	return filteredInputs
}
