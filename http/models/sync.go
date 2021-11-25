package models

type SyncOptions struct {
	PipelineSourceBranches string `url:"pipelineSourceBranches,omitempty"` // Can be a csv list
	PipelineSourceId       int    `url:"pipelineSourceId,omitempty"`
}

type SyncStatus struct {
	Id                   int        `json:"id"`
	PipelineSourceBranch string     `json:"pipelineSourceBranch"`
	PipelineSourceId     int        `json:"pipelineSourceId"`
	IsSyncing            bool       `json:"isSyncing"`
	LastSyncStatusCode   StatusCode `json:"lastSyncStatusCode"`
	LastSyncStartedAt    string     `json:"lastSyncStartedAt"`
	LastSyncEndedAt      string     `json:"lastSyncEndedAt"`
	LastSyncLogs         string     `json:"lastSyncLogs"`
	SyncUpdatedAt        string     `json:"updatedAt"`
}

type SyncStatusResponse struct {
	SyncStatuses []SyncStatus
}