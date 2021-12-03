package command

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/utils"
	"github.com/pkg/errors"
)

func New001ValidateGitState() *_001 {
	return &_001{}
}

type _001 struct{}

func (c *_001) ResolveState(ctx context.Context, state *PipedCommandState) (Status, error) {
	dirConfig := ctx.Value(utils.DirConfigCtxKey).(*utils.DirConfig)
	state.PipelinesSourceId = dirConfig.PipelinesSourceId

	branchName, err := utils.GetCurrentBranchName()
	if err != nil {
		return Status{}, errors.Wrap(err, "failed resolving current git branch")
	}

	localCommitHash, err := utils.GetCommitHash(branchName, false)
	if err != nil {
		return Status{}, err
	}
	remoteCommitHash, err := utils.GetCommitHash(branchName, true)
	if err != nil {
		return Status{}, err
	}

	if localCommitHash != remoteCommitHash {
		return Status{}, fmt.Errorf("local commit hash is different than remote, push your changes before triggering a build")
	}

	state.GitBranch = branchName
	state.HeadCommitSha = remoteCommitHash
	//state.HeadCommitSha = "b8cb635bf49ce48e6de66455b58bd374f6c84c65" //TODO only for tests

	return Status{
		Message: "git details:\ncurrent branch: %s\nlocal commit hash:  %s\nremote commit hash: %s",
		Type:    Done,
	}, nil
}

func (c *_001) TriggerStateChange(ctx context.Context, state *PipedCommandState) error {
	// do nothing
	return nil
}
