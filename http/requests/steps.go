package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	stepsUrl = "/steps"
)

func GetSteps(client http.PipelineHttpClient, options models.GetStepsOptions) (*models.StepsResponse, error) {
	body, err := client.SendGet(stepsUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.StepsResponse{}
	err = json.Unmarshal(body, &res.Steps)
	return res, err
}
