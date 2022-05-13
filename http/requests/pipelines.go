package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	pipelinesLookupUrl = "/pipelines"
)

type GetPipelinesOptions struct {
	SortBy                 string `url:"sortBy,omitempty"`
	SortOrder              int    `url:"sortOrder,omitempty"`
	FilterBy               string `url:"filterBy,omitempty"`
	Light                  bool   `url:"light,omitempty"`
	Limit                  int    `url:"limit,omitempty"`
	PipesNames             string `url:"names,omitempty"`
	PipelineSourceBranches string `url:"pipelineSourceBranches,omitempty"`
}

type Pipeline struct {
	Name                 string `json:"name"`
	PipelineId           int    `json:"id"`
	LatestRunId          int    `json:"latestRunId"`
	ProjectId            int    `json:"projectId"`
	PipelineSourceId     int    `json:"pipelineSourceId"`
	PipelineSourceBranch string `json:"pipelineSourceBranch"`
	LatestCompletedRunId int    `json:"latestCompletedRunId"`
}

type PipelinesLookupResponse struct {
	Pipelines []Pipeline
}

func GetPipelines(client http.PipelineHttpClient, options GetPipelinesOptions) (*PipelinesLookupResponse, error) {
	body, err := client.SendGet(pipelinesLookupUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &PipelinesLookupResponse{}
	err = json.Unmarshal(body, &res.Pipelines)
	return res, err
}
