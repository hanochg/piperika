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

func GetPipelinesSteps(client http.PipelineHttpClient, options models.GetPipelinesStepsOptions) (*models.PipelinesStepsResponse, error) {
	body, err := client.SendGet(stepUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.PipelinesStepsResponse{}
	err = json.Unmarshal(body, &res.Steps)
	return res, err
}

func TriggerPipelinesStep(client http.PipelineHttpClient, stepId int) {
	// Pipelines sometimes trigger the run but returns an unexpected response
	// So, we ignore the API error as it's not relevant to our use-case
	_, _ = client.SendPost(fmt.Sprintf("%s/%d/trigger", stepUrl, stepId), http.ClientOptions{Query: nil}, nil)
}
