package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	stepConnectionUrl = "/pipelineStepConnections"
)

func GetStepConnections(client http.PipelineHttpClient, options models.GetStepConnectionsOptions) (*models.StepConnectionsResponse, error) {
	body, err := client.SendGet(stepConnectionUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.StepConnectionsResponse{}
	err = json.Unmarshal(body, &res.StepConnections)
	return res, err
}
