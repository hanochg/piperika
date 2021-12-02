package command

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
)

func New002PipelinesSourcesBranchSync() *_002 {
	return &_002{}
}

type _002 struct{}

func (c *_002) ResolveState(ctx context.Context, state *PipedCommandState) (Status, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	syncStatusResp, err := requests.GetSyncStatus(httpClient, models.SyncOptions{
		PipelineSourceBranches: state.GitBranch,
		PipelineSourceId:       state.PipelinesSourceId,
		Light:                  true,
	})
	if err != nil {
		return Status{}, err
	}

	if len(syncStatusResp.SyncStatuses) == 0 {
		return Status{}, backoff.Permanent(fmt.Errorf("could not find pipes for branch %s, triggering a sync", state.GitBranch))
	}

	syncStatus := syncStatusResp.SyncStatuses[0]
	if !syncStatus.IsSyncing && syncStatus.LastSyncStatusCode != models.Success {
		return Status{}, backoff.Permanent(fmt.Errorf("sync status is complete but sync failed, triggering sync again"))
	}

	resVersions, err := requests.GetResourceVersions(httpClient, models.GetResourcesOptions{
		PipelineSourceIds:  strconv.Itoa(state.PipelinesSourceId),
		ResourceVersionIds: strconv.Itoa(syncStatus.ResourceVersionId),
	})
	if err != nil {
		return Status{}, err
	}

	if len(resVersions.Resources) == 0 {
		return Status{}, backoff.Permanent(fmt.Errorf("invalid resource version id %d for branch %s", syncStatus.ResourceVersionId, state.GitBranch))
	}

	if resVersions.Resources[0].ContentPropertyBag.CommitSha != state.HeadCommitSha {
		return Status{}, backoff.Permanent(fmt.Errorf("pipelines resource has different commit hash than the remote git commit hash, triggering a sync"))
	}

	if syncStatus.IsSyncing {
		return Status{}, fmt.Errorf("pipelines is still syncing your branch to last commit hash")
	}

	return Status{
		Message: "pipelines source synced",
	}, nil
}

func (c *_002) TriggerStateChange(ctx context.Context, state *PipedCommandState) error {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	_, err := requests.SyncSource(httpClient, models.SyncSourcesOptions{
		Branch:           state.GitBranch,
		ShouldSync:       true,
		PipelineSourceId: state.PipelinesSourceId,
	})
	if err != nil {
		return err
	}
	return nil
}
