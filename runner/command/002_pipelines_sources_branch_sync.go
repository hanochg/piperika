package command

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
)

func New002PipelinesSourcesBranchSync() *_002 {
	return &_002{}
}

type _002 struct{}

func (c *_002) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	syncStatusResp, err := requests.GetSyncStatus(httpClient, requests.SyncOptions{
		PipelineSourceBranches: state.GitBranch,
		PipelineSourceId:       state.PipelinesSourceId,
		Light:                  true,
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline sync data: %v", err),
		}
	}

	if len(syncStatusResp.SyncStatuses) == 0 {
		return Status{
			Type:            Failed,
			PipelinesStatus: "triggering sync",
			Message:         "could not find any pipeline sync data for the branch",
		}
	}

	syncStatus := syncStatusResp.SyncStatuses[0]
	if !syncStatus.IsSyncing && syncStatus.LastSyncStatusCode != http.Success {
		return Status{
			Type:            Failed,
			PipelinesStatus: "triggering sync",
			Message:         "pipeline sync for the branch previously failed, triggering new sync",
		}
	}

	resVersions, err := requests.GetResourceVersions(httpClient, requests.GetResourcesOptions{
		PipelineSourceIds:  strconv.Itoa(state.PipelinesSourceId),
		ResourceVersionIds: strconv.Itoa(syncStatus.ResourceVersionId),
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("failed fetching pipeline resources data: %v", err),
		}
	}

	if len(resVersions.Resources) == 0 {
		return Status{
			Type:            InProgress,
			PipelinesStatus: "Waiting resources",
			Message:         fmt.Sprintf("no resources for version id '%d' for branch '%s'", syncStatus.ResourceVersionId, state.GitBranch),
		}
	}

	if resVersions.Resources[0].ContentPropertyBag.CommitSha != state.HeadCommitSha {
		return Status{
			Type:            Failed,
			PipelinesStatus: "triggering sync",
			Message:         "pipelines resource has different commit hash than the remote git commit hash",
		}
	}

	if syncStatus.IsSyncing {
		return Status{
			Type:            InProgress,
			PipelinesStatus: syncStatus.LastSyncStatusCode.StatusCodeName(),
			Message:         "Pipelines is still syncing your branch to last commit hash",
		}
	}

	return Status{
		Message: fmt.Sprintf("Pipelines branch %s is synced to last commit hash", state.GitBranch),
		Type:    Done,
	}
}

func (c *_002) TriggerOnFail(ctx context.Context, state *PipedCommandState) error {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	_, err := requests.SyncSource(httpClient, requests.SyncSourcesOptions{
		Branch:           state.GitBranch,
		ShouldSync:       true,
		PipelineSourceId: state.PipelinesSourceId,
	})
	return err
}
