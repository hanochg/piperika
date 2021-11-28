package requests

import (
	"encoding/json"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	sourcesUrl = "/pipelineSources"
)

func GetSource(client http.PipelineHttpClient, options models.GetSourcesOptions) (*models.SourcesResponse, error) {
	body, err := client.SendGet(sourcesUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.SourcesResponse{}
	err = json.Unmarshal(body, &res.Sources)
	return res, err
}

func SyncSource(client http.PipelineHttpClient, options models.SyncSourcesOptions) (*models.SourcesResponse, error) {
	body, err := client.SendGet(sourcesUrl+fmt.Sprintf("/%d", options.PipelineSourceId), http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.SourcesResponse{Sources: []models.Source{{}}}
	err = json.Unmarshal(body, &res.Sources[0])
	return res, err
}
