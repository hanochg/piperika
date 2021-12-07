package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	syncUrl = "/pipelineSyncStatuses"
)

type SyncOptions struct {
	PipelineSourceBranches string `url:"pipelineSourceBranches,omitempty"` // Can be a csv list
	PipelineSourceId       int    `url:"pipelineSourceId,omitempty"`
	Light                  bool   `url:"light,omitempty"`
}

type SyncStatus struct {
	Id                   int             `json:"id"`
	PipelineSourceBranch string          `json:"pipelineSourceBranch"`
	PipelineSourceId     int             `json:"pipelineSourceId"`
	IsSyncing            bool            `json:"isSyncing"`
	LastSyncStatusCode   http.StatusCode `json:"lastSyncStatusCode"`
	LastSyncStartedAt    string          `json:"lastSyncStartedAt"`
	LastSyncEndedAt      string          `json:"lastSyncEndedAt"`
	LastSyncLogs         string          `json:"lastSyncLogs"`
	SyncUpdatedAt        string          `json:"updatedAt"`
	ResourceVersionId    int             `json:"triggeredByResourceVersionId"`
}

type SyncStatusResponse struct {
	SyncStatuses []SyncStatus
}

func GetSyncStatus(client http.PipelineHttpClient, options SyncOptions) (*SyncStatusResponse, error) {
	body, err := client.SendGet(syncUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &SyncStatusResponse{}
	err = json.Unmarshal(body, &res.SyncStatuses)
	return res, err
}
