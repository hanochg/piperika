package runs

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/runner/datastruct"
	"github.com/hanochg/piperika/utils"
	"github.com/pkg/errors"
	"strconv"
)

/*
	Run Description
	---------------
	- Validate the local and remote git commits are similar (developer pushed its code).
	- Check if the branch exists on Pipelines resources
	- Check if the last branch sync was successful
	- Check if the synced branch is synced to the right commit sha
*/

func (_ _01) Init(ctx context.Context, state *datastruct.PipedCommandState) (string, error) {
	branchName, err := utils.GetCurrentBranchName()
	if err != nil {
		return "", errors.Wrap(err, "failed resolving current git branch")
	}

	localCommitHash, err := utils.GetCommitHash(branchName, false)
	if err != nil {
		return "", err
	}
	remoteCommitHash, err := utils.GetCommitHash(branchName, true)
	if err != nil {
		return "", err
	}

	if localCommitHash != remoteCommitHash {
		return "", fmt.Errorf("local commit hash is different than remote, push your changes before triggering a build")
	}

	state.GitBranch = branchName
	state.HeadCommitSha = remoteCommitHash
	//state.HeadCommitSha = "b8cb635bf49ce48e6de66455b58bd374f6c84c65" //TODO only for tests

	return "git details:\ncurrent branch: %s\nlocal commit hash:  %s\nremote commit hash: %s", nil
}

func (_ _01) Tick(ctx context.Context, state *datastruct.PipedCommandState) (*datastruct.RunStatus, error) {
	httpClient := ctx.Value("client").(http.PipelineHttpClient)
	dirConfig := ctx.Value("dirConfig").(*utils.DirConfig)
	state.ShouldTriggerPipelinesSync = true
	syncStatusResp, err := requests.GetSyncStatus(httpClient, models.SyncOptions{
		PipelineSourceBranches: state.GitBranch,
		PipelineSourceId:       dirConfig.PipelinesSourceId,
		Light:                  true,
	})
	if err != nil {
		return nil, err
	}

	if len(syncStatusResp.SyncStatuses) == 0 {
		return &datastruct.RunStatus{
			Message: fmt.Sprintf("Couldn't find pipes for branch %s, retrying", state.GitBranch),
			Status:  "waiting for pipeline",
			Done:    false,
		}, nil
	}

	syncStatus := syncStatusResp.SyncStatuses[0]
	if !syncStatus.IsSyncing && syncStatus.LastSyncStatusCode != models.Success {
		return &datastruct.RunStatus{
			Message: fmt.Sprintf("Sync status is complete but sync failed, triggering sync again"),
			Status:  "sync failed",
			Done:    true,
		}, nil
	}

	resVersions, err := requests.GetResourceVersions(httpClient, models.GetResourcesOptions{
		PipelineSourceIds:  strconv.Itoa(utils.ArtifactoryPipelinesSourceId),
		ResourceVersionIds: strconv.Itoa(syncStatus.ResourceVersionId),
	})
	if err != nil {
		return nil, err
	}

	if len(resVersions.Resources) == 0 {
		return nil, fmt.Errorf("invalid resource number %d for branch %s", syncStatus.ResourceVersionId, state.GitBranch)
	}

	if resVersions.Resources[0].ContentPropertyBag.CommitSha != state.HeadCommitSha {
		fmt.Println("Pipelines resource has different commit hash than the remote git commit hash, syncing...")
		return &datastruct.RunStatus{
			Message: fmt.Sprintf("different commit hashes"),
			Done:    true,
		}, nil
	}

	state.ShouldTriggerPipelinesSync = false
	return &datastruct.RunStatus{
		Done: true,
	}, nil
}

func (_ _01) OnComplete(ctx context.Context, state *datastruct.PipedCommandState, status *datastruct.RunStatus) (string, error) {
	httpClient := ctx.Value("client").(http.PipelineHttpClient)
	if state.ShouldTriggerPipelinesSync {
		_, err := requests.SyncSource(httpClient, models.SyncSourcesOptions{
			Branch:           state.GitBranch,
			ShouldSync:       true,
			PipelineSourceId: utils.ArtifactoryPipelinesSourceId,
		})
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

type _01 struct {
}

func New001GetPipeSourceBranch() datastruct.Runner {
	return _01{}
}
