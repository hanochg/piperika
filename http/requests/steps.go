package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	stepsUrl = "/steps"
)

type GetStepsOptions struct {
	RunIds string `url:"runIds,omitempty"`
	Limit  int    `url:"limit,omitempty"`
}

type Step struct {
	Id         int             `json:"id"`
	Name       string          `json:"name"`
	StatusCode http.StatusCode `json:"statusCode"`
}

type StepsResponse struct {
	Steps []Step
}

func GetSteps(client http.PipelineHttpClient, options GetStepsOptions) (*StepsResponse, error) {
	body, err := client.SendGet(stepsUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &StepsResponse{}
	err = json.Unmarshal(body, &res.Steps)
	return res, err
}
