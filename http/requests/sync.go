package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	syncUrl = "/pipelineSyncStatuses"
)

func GetSyncStatus(client http.PipelineHttpClient, options models.SyncOptions) (*models.SyncStatusResponse, error) {
	body, err := client.SendGet(syncUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.SyncStatusResponse{}
	err = json.Unmarshal(body, res.SyncStatuses)
	return res, err
}
