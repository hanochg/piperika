package datastruct

import (
	"context"
	"time"
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
	PipelinesSyncDate time.Time
	PipelineId        int
	RunId             int
}

type PipedCommand interface {
	Run(ctx context.Context, state *PipedCommandState) error
}
