package requests

import (
	"encoding/json"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	stepUrl = "/pipelineSteps"
)

func GetPipelinesSteps(client http.PipelineHttpClient, options models.GetStepsOptions) (*models.StepsResponse, error) {
	body, err := client.SendGet(stepUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.StepsResponse{}
	err = json.Unmarshal(body, &res.Steps)
	return res, err
}

func TriggerPipelinesStep(client http.PipelineHttpClient, options models.GetStepsOptions, stepId int) (*models.StepsResponse, error) {
	body, err := client.SendPost(fmt.Sprintf("%s/%d/trigger", stepUrl, stepId), http.ClientOptions{Query: options}, nil)
	if err != nil {
		return nil, err
	}
	res := &models.StepsResponse{}
	err = json.Unmarshal(body, &res.Steps)
	return res, err
}
