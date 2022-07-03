package commands

import (
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"time"
)

var finalStatus = []http.StatusCode{http.Success, http.Error, http.Canceled, http.Failure}

func GetWaitSyncCommand() components.Command {
	return components.Command{
		Name:        "wait-sync",
		Description: "Wait for sync to complete",
		Aliases:     []string{"ws"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		Action:      waitSync,
	}
}

func waitSync(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}

	config, err := utils.GetConfigurations()
	if err != nil {
		return err
	}

	branchName, err := utils.GetCurrentBranchName(c)
	if err != nil {
		return err
	}

	timeout := time.Minute * 2  // TODO: make as flag
	interval := time.Second * 3 // TODO: make as flag

	timeoutTimer := time.NewTimer(timeout)
	lastStatus, err := getSyncStatus(client, branchName, config)
	for keepCheckSync(lastStatus) {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf("timeout")
		case <-time.After(interval):
			lastStatus, err = getSyncStatus(client, branchName, config)
			if err != nil {
				return err
			}
		}
	}

	if *lastStatus != http.Success {
		return fmt.Errorf("unexpected status %s", lastStatus.String())
	}
	return nil
}

func getSyncStatus(client http.PipelineHttpClient, branchName string, config *utils.Configurations) (*http.StatusCode, error) {
	status, err := requests.GetSyncStatus(client, requests.SyncOptions{
		PipelineSourceBranches: branchName,
		PipelineSourceId:       config.PipelinesSourceId,
		Light:                  true,
	})
	if err != nil {
		return nil, err
	}

	if len(status.SyncStatuses) > 1 {
		return nil, fmt.Errorf("more then 1 syncs: %v", status.SyncStatuses)
	} else if len(status.SyncStatuses) == 0 {
		return nil, nil
	}

	return &status.SyncStatuses[0].LastSyncStatusCode, nil
}

func keepCheckSync(lastStatus *http.StatusCode) bool {
	if lastStatus == nil {
		return true
	}
	for _, checkStatus := range finalStatus {
		if checkStatus == *lastStatus {
			return false
		}
	}

	return true
}
