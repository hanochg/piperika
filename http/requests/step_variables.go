package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	stepVariablesUrl = "/stepVariables"
)

type GetStepVariablesOptions struct {
	StepIds int `url:"stepIds,omitempty"`
}

type StepVariable struct {
	Id         int       `json:"id"`
	ProjectId  int       `json:"projectId"`
	PipelineId int       `json:"pipelineId"`
	RunId      int       `json:"runId"`
	StepId     int       `json:"stepId"`
	StepletId  int       `json:"stepletId"`
	Variables  Variables `json:"variables"`
}

type Variables struct {
	RunVariable []RunVariable `json:"runVariables"`
}

type RunVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StepVariablesResponse struct {
	Variables []StepVariable
}

func GetStepVariables(client http.PipelineHttpClient, options GetStepVariablesOptions) (*StepVariablesResponse, error) {
	body, err := client.SendGet(stepVariablesUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &StepVariablesResponse{}
	err = json.Unmarshal(body, &res.Variables)
	return res, err
}
