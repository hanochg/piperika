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
			Message: fmt.Sprintf("failed resolving current git branch: %v", err),
		}
	}

	localCommitHash, err := utils.GetCommitHash(branchName, false)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("failed resolving local git commit hash: %v", err),
		}
	}
	remoteCommitHash, err := utils.GetCommitHash(branchName, true)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("failed resolving remote git commit hash: %v", err),
		}
	}

	if localCommitHash != remoteCommitHash {
		return Status{
			Type:    Unrecoverable,
			Message: "local commit hash is different than remote, push your changes before triggering a build",
		}
	}

	state.GitBranch = branchName
	state.HeadCommitSha = remoteCommitHash
	//state.HeadCommitSha = "b8cb635bf49ce48e6de66455b58bd374f6c84c65" //TODO only for tests

	return Status{
		Message: fmt.Sprintf("git details:\ncurrent branch: %s\nlocal commit hash:  %s\nremote commit hash: %s",
			branchName, localCommitHash, remoteCommitHash),
		Type: Done,
	}
}

func (c *_001) TriggerStateChange(ctx context.Context, state *PipedCommandState) Status {
	return Status{
		Type:    Unrecoverable,
		Message: "Timed out",
	}
}
