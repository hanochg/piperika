package requests

import (
	"encoding/json"
	"fmt"
	"github.com/hanochg/piperika/http"
)

const (
	sourcesUrl = "/pipelineSources"
)

type GetSourcesOptions struct {
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"` // Can be a csv list
}

type SyncSourcesOptions struct {
	Branch           string `url:"branch,omitempty"`
	ShouldSync       bool   `url:"sync,omitempty"`
	PipelineSourceId int
}

type Source struct {
	Id                 int             `json:"id"`
	RepositoryFullName string          `json:"repositoryFullName"`
	LastSyncStatusCode http.StatusCode `json:"lastSyncStatusCode"`
	IsSyncing          bool            `json:"isSyncing"`
	LastSyncStartedAt  string          `json:"lastSyncStartedAt"`
	LastSyncEndedAt    string          `json:"lastSyncEndedAt"`
	LastSyncLogs       string          `json:"lastSyncLogs"`
	SyncUpdatedAt      string          `json:"updatedAt"`
}

type SourcesResponse struct {
	Sources []Source
}

func SyncSource(client http.PipelineHttpClient, options SyncSourcesOptions) (*SourcesResponse, error) {
	body, err := client.SendGet(sourcesUrl+fmt.Sprintf("/%d", options.PipelineSourceId), http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &SourcesResponse{Sources: []Source{{}}}
	err = json.Unmarshal(body, &res.Sources[0])
	return res, err
}
