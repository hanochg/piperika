package models

type GetRunsOptions struct {
	PipelineIds string `url:"pipelineIds,omitempty"`
	Limit       int    `url:"limit,omitempty"`
	Light       bool   `url:"light,omitempty"`
	StatusCodes string `url:"statusCodes,omitempty"` //Optional. Comma separated string of status codes
	SortBy      string `url:"sortBy,omitempty"`      //Optional. Comma separated list of sort attributes
	SortOrder   int    `url:"sortOrder,omitempty"`   //Optional. 1 for ascending and -1 for descending based on sortBy
	RunNumbers  string `url:"sortOrder,omitempty"`
	RunIds      string `url:"runIds,omitempty"`
}

type Run struct {
	RunId            int        `json:"id"`
	PipelineId       int        `json:"pipelineId"`
	PipelineSourceId int        `json:"pipelineSourceId"`
	RunNumber        int        `json:"runNumber"`
	StatusCode       StatusCode `json:"statusCode"`
	StartedAt        string     `json:"startedAt"`
	EndedAt          string     `json:"endedAt"`
}

type RunsResponse struct {
	Runs []Run
}
