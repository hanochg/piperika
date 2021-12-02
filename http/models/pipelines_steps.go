package models

type GetPipelinesStepsOptions struct {
	PipelineIds       string `url:"pipelineIds,omitempty"`
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"`
	Names             string `url:"names,omitempty"`
}

type PipelinesSteps struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PipelinesStepsResponse struct {
	Steps []PipelinesSteps
}
