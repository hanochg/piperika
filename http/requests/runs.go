package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	runsUrl = "/runs"
)

type GetRunsOptions struct {
	PipelineIds string `url:"pipelineIds,omitempty"`
	Limit       int    `url:"limit,omitempty"`
	Light       bool   `url:"light,omitempty"`
	StatusCodes string `url:"statusCodes,omitempty"` //Optional. Comma separated string of status codes
	SortBy      string `url:"sortBy,omitempty"`      //Optional. Comma separated list of sort attributes
	SortOrder   int    `url:"sortOrder,omitempty"`   //Optional. 1 for ascending and -1 for descending based on sortBy
	RunNumbers  string `url:"sortOrder,omitempty"`
	RunIds      string `url:"runIds,omitempty"`
	ProjectIds  string `url:"projectIds,omitempty"`
}

type Run struct {
	RunId             int               `json:"id"`
	PipelineId        int               `json:"pipelineId"`
	PipelineSourceId  int               `json:"pipelineSourceId"`
	RunNumber         int               `json:"runNumber"`
	StatusCode        http.StatusCode   `json:"statusCode"`
	StartedAt         string            `json:"startedAt"`
	EndedAt           string            `json:"endedAt"`
	TriggeredAt       string            `json:"createdAt"`
	StaticPropertyBag StaticPropertyBag `json:"staticPropertyBag"`
	ProjectId         int               `json:"projectId"`
}

type StaticPropertyBag struct {
	TriggeredByUserName     string `json:"triggeredByUserName"`
	TriggeredByResourceName string `json:"triggeredByResourceName"`
	SignedPipelinesEnabled  bool   `json:"signedPipelinesEnabled"`
}

type RunsResponse struct {
	Runs []Run
}

func GetRuns(client http.PipelineHttpClient, options GetRunsOptions) (*RunsResponse, error) {
	body, err := client.SendGet(runsUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &RunsResponse{}
	err = json.Unmarshal(body, &res.Runs)
	return res, err
}
