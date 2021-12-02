package models

type GetPipelinesOptions struct {
	SortBy     string `url:"sortBy,omitempty"`
	FilterBy   string `url:"filterBy,omitempty"`
	Light      bool   `url:"light,omitempty"`
	Limit      int    `url:"limit,omitempty"`
	PipesNames string `url:"names,omitempty"`
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
