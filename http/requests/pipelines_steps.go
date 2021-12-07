package requests

import (
	"encoding/json"
	"fmt"
	"github.com/hanochg/piperika/http"
)

const (
	stepUrl = "/pipelineSteps"
)

type GetPipelinesStepsOptions struct {
	PipelineIds       string `url:"pipelineIds,omitempty"`
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"`
	Names             string `url:"names,omitempty"`
}

type PipelinesSteps struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PipelinesStepsResponse struct {
	Steps []PipelinesSteps
}

func GetPipelinesSteps(client http.PipelineHttpClient, options GetPipelinesStepsOptions) (*PipelinesStepsResponse, error) {
	body, err := client.SendGet(stepUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &PipelinesStepsResponse{}
	err = json.Unmarshal(body, &res.Steps)
	return res, err
}

func TriggerPipelinesStep(client http.PipelineHttpClient, stepId int) error {
	_, err := client.SendPost(fmt.Sprintf("%s/%d/trigger", stepUrl, stepId), http.ClientOptions{Query: nil}, nil)
	return err
}
