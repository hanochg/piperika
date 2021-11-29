package runner

import "context"

type pipedCommandState struct {
	// git details
	GitBranch        string
	LocalCommitHash  string
	RemoteCommitHash string

	// trigger pipelines
	ShouldTriggerPipelinesSync bool
}

type pipedCommand interface {
	Run(ctx context.Context, state *pipedCommandState) error
}
