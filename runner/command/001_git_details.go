package command

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/utils"
	"github.com/pkg/errors"
)

func NewGitDetailsCommand() PipedCommand {
	const opName = "git details"
	initialBackoff := &backoff.StopBackOff{}

	return NewRetryingPipedCommand(opName, initialBackoff, func(ctx context.Context, state *PipedCommandState) error {
		branchName, err := utils.GetCurrentBranchName()
		if err != nil {
			return backoff.Permanent(errors.Wrap(err, "failed resolving current git branch"))
		}

		// TODO
		localCommitHash := ""
		remoteCommitHash := ""
		if localCommitHash != remoteCommitHash {
			return backoff.Permanent(fmt.Errorf("local commit hash is different than remote, push your changes before triggering a build"))
		}

		logInfo(opName, fmt.Sprintf("git details:\ncurrent branch: %s\nlocal commit hash:  %s\nremote commit hash: %s",
			branchName, localCommitHash, remoteCommitHash))

		state.GitBranch = branchName
		state.LocalCommitHash = localCommitHash
		state.RemoteCommitHash = remoteCommitHash
		return nil
	})
}
