package command

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/utils"
)

func New001ValidateGitState() *_001 {
	return &_001{}
}

type _001 struct{}

func (c *_001) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	dirConfig := ctx.Value(utils.DirConfigCtxKey).(*utils.DirConfig)
	state.PipelinesSourceId = dirConfig.PipelinesSourceId

	branchName, err := utils.GetCurrentBranchName()
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed resolving current git branch: %v", err),
		}
	}

	localCommitHash, err := utils.GetCommitHash(branchName, false)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed resolving local git commit hash: %v", err),
		}
	}
	remoteCommitHash, err := utils.GetCommitHash(branchName, true)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed resolving remote git commit hash: %v", err),
		}
	}

	if localCommitHash != remoteCommitHash {
		return Status{
			Type:    Unrecoverable,
			Message: "Local commit hash is different than remote, push your changes before triggering a build",
		}
	}

	state.GitBranch = branchName
	state.HeadCommitSha = remoteCommitHash

	return Status{
		Type: Done,
	}
}

func (c *_001) TriggerOnFail(_ context.Context, _ *PipedCommandState) error {
	return fmt.Errorf("timed out")
}
