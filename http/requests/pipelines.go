package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	pipelinesLookupUrl = "/pipelines"
)

func GetPipelines(client http.PipelineHttpClient, options models.PipelinesLookupOptions) (*models.PipelinesLookupResponse, error) {
	body, err := client.SendGet(pipelinesLookupUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.PipelinesLookupResponse{}
	err = json.Unmarshal(body, &res.Pipelines)
	return res, err

}
