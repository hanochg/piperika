package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
)

const (
	stepsTestReportsUrl = "/stepTestReports"
)

func GetStepsTestReports(client http.PipelineHttpClient, options models.StepsTestReportsOptions) (*models.GetStepsTestReportResponse, error) {
	body, err := client.SendGet(stepsTestReportsUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &models.GetStepsTestReportResponse{}
	err = json.Unmarshal(body, &res.TestReports)
	return res, err
}
