package runner

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/utils"
	"github.com/pkg/errors"
)

func (_ _01) init(ctx context.Context, state *pipedCommandState) (string, error) {
	branchName, err := utils.GetCurrentBranchName()
	if err != nil {
		return "", errors.Wrap(err, "failed resolving current git branch")
	}

	// TODO
	localCommitHash := ""
	remoteCommitHash := ""
	if localCommitHash != remoteCommitHash {
		return "", fmt.Errorf("local commit hash is different than remote, push your changes before triggering a build")
	}

	state.GitBranch = branchName
	state.LocalCommitHash = localCommitHash
	state.RemoteCommitHash = remoteCommitHash

	return "git details:\ncurrent branch: %s\nlocal commit hash:  %s\nremote commit hash: %s", nil
}

func (_ _01) tick(ctx context.Context, state *pipedCommandState) (*tickStatus, error) {
	// TODO get pipeline source branch and wait for existence
	return nil, nil
}

func (_ _01) timedOutOperation(ctx context.Context, state *pipedCommandState, status *tickStatus) (string, error) {
	// TODO if not exists, fetch by branch
	// client := ctx.Value("client").(http.PipelineHttpClient) // Getting a client, consider use helper for ease or dedicated ctx
	return "", nil
}

type _01 struct {
}

func new001GetPipeSourceBranch() _01 {
	return _01{}
}
