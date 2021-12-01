package steps

import "github.com/hanochg/piperika/http/models"

type GetStepsOptions struct {
	PipelineIds       string            `url:"pipelineIds,omitempty"`
	StatusCode        models.StatusCode `url:"statusCode,omitempty"`
	Limit             int               `url:"limit,omitempty"`
	StepIds           string            `url:"stepIds,omitempty"`
	PipelineSourceIds string            `url:"pipelineSourceIds,omitempty"`
}

type Step struct {
	Name             string            `json:"name"`
	PipelineId       int               `json:"pipelineId"`
	PipelineSourceId int               `json:"pipelineSourceId"`
	PipelineStepId   int               `json:"pipelineStepId"`
	RunId            int               `json:"runId"`
	StatusCode       models.StatusCode `json:"statusCode"`
}

type GetStepsResponse struct {
	Steps []Step
}
