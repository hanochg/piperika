package models

type GetSourcesOptions struct {
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"` // Can be a csv list
}

type SyncSourcesOptions struct {
	Branch           string `url:"branch,omitempty"`
	ShouldSync       bool   `url:"sync,omitempty"`
	PipelineSourceId int
}

type Source struct {
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
	Sources []Source
}
