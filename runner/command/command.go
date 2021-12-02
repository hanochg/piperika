package command

import "context"

type PipedCommandState struct {
	// git details
	GitBranch     string
	HeadCommitSha string

	// triggering
	ShouldTriggerPipelinesSync bool
	ShouldTriggerRun           bool

	// pipelines details
	PipelinesSourceId int
	PipelineId        int
	RunId             int
	RunNumber         int
	RunStepIdsCsv     string
}

type Status struct {
	PipelinesStatus string
	Message         string
}

type Command interface {
	ResolveState(ctx context.Context, state *PipedCommandState) (status Status, err error)
	TriggerStateChange(ctx context.Context, state *PipedCommandState) error
}
