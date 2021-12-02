package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	runResourcesUrl = "/runResourceVersions"
)

func GetRunResourceVersions(client http.PipelineHttpClient, options models.GetRunResourcesOptions) (*models.RunResourcesResponse, error) {
	body, err := client.SendGet(runResourcesUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.RunResourcesResponse{}
	err = json.Unmarshal(body, &res.Resources)
	return res, err
}
