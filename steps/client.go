package steps

import (
	"github.com/hanochg/piperika/utils"
)

const (
	stepsUrl = "steps"
)

type GetStepsOptions struct {
	pipelineIds string "url:pipelineIds"
}

func getSteps(client utils.PipelineHttpClient, options GetStepsOptions) (interface{}, error) {
	return client.SendGet(stepsUrl, utils.ClientOptions{Query: options})

}
