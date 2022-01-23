package command

import "context"

type PipedCommandState struct {
	// git details
	HeadCommitSha string

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
	DynamicMessage  func(size int)
	Type            StatusType
	Link            string
}

type StatusType string

const (
	InProgress    StatusType = "in-progress"
	Done          StatusType = "done"
	Failed        StatusType = "failed"
	Unrecoverable StatusType = "unrecoverable"
)

type Command interface {
	ResolveState(ctx context.Context, state *PipedCommandState) Status
	TriggerOnFail(ctx context.Context, state *PipedCommandState) error
}
