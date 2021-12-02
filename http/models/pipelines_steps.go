package models

type GetStepsOptions struct {
	PipelineIds       string `url:"pipelineIds,omitempty"`
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"`
	Names             string `url:"names,omitempty"`
}

type Steps struct {
	Id   int `json:"id"`
	Name int `json:"name"`
}

type StepsResponse struct {
	Steps []Steps
}
