package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	sourcesUrl = "/pipelineSources"
)

func GetSource(client http.PipelineHttpClient, options models.SourcesOptions) (*models.SourcesResponse, error) {
	body, err := client.SendGet(sourcesUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.SourcesResponse{}
	err = json.Unmarshal(body, &res.Sources)
	return res, err
}
