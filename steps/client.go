package steps

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	stepsUrl = "/steps"
)

func GetSteps(client http.PipelineHttpClient, options GetStepsOptions) (*GetStepsResponse, error) {
	body, err := client.SendGet(stepsUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &GetStepsResponse{}
	err = json.Unmarshal(body, &res.Steps)
	return res, err
}
