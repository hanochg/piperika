package command

import "context"

type PipedCommandState struct {
	// git details
	GitBranch        string
	LocalCommitHash  string
	RemoteCommitHash string

	// trigger pipelines
	ShouldTriggerPipelinesSync bool
}

type CommandResolution struct {
	Error error
}

type PipedCommand interface {
	Run(ctx context.Context, state *PipedCommandState) CommandResolution
}
