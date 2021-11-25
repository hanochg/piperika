package models

type SourcesOptions struct {
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"` // Can be a csv list
}

type Sources struct {
	Id                 int        `json:"id"`
	RepositoryFullName string     `json:"repositoryFullName"`
	LastSyncStatusCode StatusCode `json:"lastSyncStatusCode"`
	IsSyncing          bool       `json:"isSyncing"`
	LastSyncStartedAt  string     `json:"lastSyncStartedAt"`
	LastSyncEndedAt    string     `json:"lastSyncEndedAt"`
	LastSyncLogs       string     `json:"lastSyncLogs"`
	SyncUpdatedAt      string     `json:"updatedAt"`
}

type SourcesResponse struct {
	Sources []Sources
}
