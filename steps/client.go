package steps

import (
	"encoding/json"
	"github.com/hanochg/piperika/utils"
)

const (
	stepsUrl = "steps"
)

type GetStepsOptions struct {
	PipelineIds string     `url:"pipelineIds"`
	StatusCode  StatusCode `url:"statusCode"`
	Limit       int        `url:"limit"`
}

type GetStepsResponse struct {
	Steps []Step `json:"steps"`
}

func getSteps(client utils.PipelineHttpClient, options GetStepsOptions) (*GetStepsResponse, error) {
	body, err := client.SendGet(stepsUrl, utils.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &GetStepsResponse{}
	err = json.Unmarshal(body, res)
	return res, err

}
