package steps

import (
	"github.com/hanochg/piperika/utils"
)

func GetRunningStepsForBranch(client utils.PipelineHttpClient, branch string) ([]string, error) {
	// TODO fetch pipeline and last run based on branch using other services
	body, err := getSteps(client, GetStepsOptions{
		StatusCode: InProgress,
	})
	if err != nil {
		return nil, err
	}
	res := make([]string, len(body.Steps))
	for i, step := range body.Steps {
		res[i] = step.Name
	}

	return res, nil
}
