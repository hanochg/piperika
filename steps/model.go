package steps

type GetStepsOptions struct {
	PipelineIds       string     `url:"pipelineIds,omitempty"`
	StatusCode        StatusCode `url:"statusCode,omitempty"`
	Limit             int        `url:"limit,omitempty"`
	StepIds           string     `url:"stepIds,omitempty"`
	PipelineSourceIds string     `url:"pipelineSourceIds,omitempty"`
}

type Step struct {
	Name             string     `json:"name"`
	PipelineId       int        `json:"pipelineId"`
	PipelineSourceId int        `json:"pipelineSourceId"`
	PipelineStepId   int        `json:"pipelineStepId"`
	RunId            int        `json:"runId"`
	StatusCode       StatusCode `json:"statusCode"`
}

type StatusCode int

const (
	Queued     StatusCode = 4000
	InProgress StatusCode = 4001
	Success    StatusCode = 4002
	Failure    StatusCode = 4003
	Error      StatusCode = 4004
	Wait       StatusCode = 4005
	Canceled   StatusCode = 4006
	Unstable   StatusCode = 4007
	Skipped    StatusCode = 4008
	TimedOut   StatusCode = 4009
	TimingOut  StatusCode = 4014
)

type GetStepsResponse struct {
	Steps []Step
}
