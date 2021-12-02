package requests

import (
	"encoding/json"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	runsUrl      = "/runs"
	cancelRunUrl = "/runs/%d/cancel"
)

func GetRuns(client http.PipelineHttpClient, options models.GetRunsOptions) (*models.RunsResponse, error) {
	body, err := client.SendGet(runsUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.RunsResponse{}
	err = json.Unmarshal(body, &res.Runs)
	return res, err
}

func CancelRun(client http.PipelineHttpClient, runId int) (*models.RunsResponse, error) {

	body, err := client.SendPost(fmt.Sprintf(cancelRunUrl, runId), http.ClientOptions{}, nil)
	if err != nil {
		return nil, err
	}
	res := &models.RunsResponse{}
	err = json.Unmarshal(body, &res.Runs)
	return res, err
}
