package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	resourcesUrl = "/resourceVersions"
)

func GetResourceVersions(client http.PipelineHttpClient, options models.GetResourcesOptions) (*models.ResourcesResponse, error) {
	body, err := client.SendGet(resourcesUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.ResourcesResponse{}
	err = json.Unmarshal(body, &res.Resources)
	return res, err
}
