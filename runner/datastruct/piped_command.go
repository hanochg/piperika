package datastruct

import (
	"context"
)

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

type PipedCommand interface {
	Run(ctx context.Context, state *PipedCommandState) error
}
