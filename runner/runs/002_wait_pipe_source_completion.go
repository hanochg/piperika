package runs

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/runner/datastruct"
	"github.com/hanochg/piperika/utils"
)

/*
	Run Description
	---------------
	- Wait for a successful sync completion
*/

func (_ _02) Init(ctx context.Context, state *datastruct.PipedCommandState) (string, error) {
	return "", nil
}

func (_ _02) Tick(ctx context.Context, state *datastruct.PipedCommandState) (*datastruct.RunStatus, error) {
	httpClient := ctx.Value("client").(http.PipelineHttpClient)
	syncStatus, err := requests.GetSyncStatus(httpClient, models.SyncOptions{
		PipelineSourceBranches: state.GitBranch,
		PipelineSourceId:       utils.ArtifactoryPipelinesSourceId,
		Light:                  false,
	})
	if err != nil {
		return nil, err
	}

	if len(syncStatus.SyncStatuses) == 0 {
		return nil, fmt.Errorf("could not fetch branch sync status, branch %s", state.GitBranch)
	}

	if syncStatus.SyncStatuses[0].IsSyncing {
		return &datastruct.RunStatus{
			Message: "still syncing",
			Done:    false,
		}, nil
	}
	return &datastruct.RunStatus{
		Message: "nothing to sync or wait for",
		Done:    true,
	}, nil
}

func (_ _02) OnComplete(ctx context.Context, state *datastruct.PipedCommandState, status *datastruct.RunStatus) (string, error) {
	httpClient := ctx.Value("client").(http.PipelineHttpClient)
	syncStatusResp, err := requests.GetSyncStatus(httpClient, models.SyncOptions{
		PipelineSourceBranches: state.GitBranch,
		PipelineSourceId:       utils.ArtifactoryPipelinesSourceId,
		Light:                  false,
	})
	if err != nil {
		return "", err
	}

	if len(syncStatusResp.SyncStatuses) == 0 {
		return "", fmt.Errorf("could not fetch branch sync status, branch %s", state.GitBranch)
	}

	if syncStatusResp.SyncStatuses[0].IsSyncing ||
		syncStatusResp.SyncStatuses[0].LastSyncStatusCode != models.Success {
		return "", fmt.Errorf("branch sync is not ready or faulty, branch %s, is syncing %t, last sync status %d",
			state.GitBranch, syncStatusResp.SyncStatuses[0].IsSyncing, syncStatusResp.SyncStatuses[0].LastSyncStatusCode)
	}

	lastSynced, err := utils.PipelinesTimeParser(syncStatusResp.SyncStatuses[0].LastSyncEndedAt)
	if err != nil {
		return "", err
	}
	state.PipelinesSyncDate = lastSynced
	return "", nil
}

type _02 struct {
}

func New002WaitPipSourceCompletion() datastruct.Runner {
	return _02{}
}
