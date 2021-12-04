package command

import (
	"context"
	"fmt"
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

func (c *_002) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	syncStatusResp, err := requests.GetSyncStatus(httpClient, models.SyncOptions{
		PipelineSourceBranches: state.GitBranch,
		PipelineSourceId:       state.PipelinesSourceId,
		Light:                  true,
	})
	if err != nil {
		return Status{
			Type:    InProgress,
			Message: fmt.Sprintf("Failed fetching pipeline sync data: %v", err),
		}
	}

	if len(syncStatusResp.SyncStatuses) == 0 {
		return Status{
			Type:    Failed,
			Message: "Could not find any pipeline sync data for the branch, triggering sync",
		}
	}

	syncStatus := syncStatusResp.SyncStatuses[0]
	if !syncStatus.IsSyncing && syncStatus.LastSyncStatusCode != models.Success {
		return Status{
			Type:    Failed,
			Message: "Pipeline sync for the branch has already run and failed, triggering new sync",
		}
	}

	resVersions, err := requests.GetResourceVersions(httpClient, models.GetResourcesOptions{
		PipelineSourceIds:  strconv.Itoa(state.PipelinesSourceId),
		ResourceVersionIds: strconv.Itoa(syncStatus.ResourceVersionId),
	})
	if err != nil {
		return Status{
			Type:    InProgress,
			Message: fmt.Sprintf("Failed fetching pipeline resources data: %v", err),
		}
	}

	if len(resVersions.Resources) == 0 {
		return Status{
			Type:    InProgress,
			Message: fmt.Sprintf("No resources for version id '%d' for branch '%s'", syncStatus.ResourceVersionId, state.GitBranch),
		}
	}

	if resVersions.Resources[0].ContentPropertyBag.CommitSha != state.HeadCommitSha {
		return Status{
			Type:    Failed,
			Message: "Pipelines resource has different commit hash than the remote git commit hash, triggering a sync",
		}
	}

	if syncStatus.IsSyncing {
		return Status{
			PipelinesStatus: fmt.Sprintf("%d", syncStatus.LastSyncStatusCode), // Map to string
			Message:         "pipelines is still syncing your branch to last commit hash",
			Type:            InProgress,
		}
	}

	return Status{
		Message: "pipelines source synced",
		Type:    Done,
	}
}

func (c *_002) TriggerStateChange(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	_, err := requests.SyncSource(httpClient, models.SyncSourcesOptions{
		Branch:           state.GitBranch,
		ShouldSync:       true,
		PipelineSourceId: state.PipelinesSourceId,
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed triggering pipeline sync: %v", err),
		}
	}
	return Status{
		Type:    Done,
		Message: "Successfully triggered pipeline sync",
	}
}
