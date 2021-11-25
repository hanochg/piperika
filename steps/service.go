package steps

import (
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
)

func GetRunningStepsForBranch(client http.PipelineHttpClient, branch string) ([]string, error) {
	// TODO fetch pipeline and last run based on branch using other services
	body, err := requests.GetSteps(client, models.GetStepsOptions{
		StatusCode:        models.InProgress, // TODO parameter
		Limit:             2,                 // TODO const
		PipelineSourceIds: "6",
	})
	if err != nil {
		return nil, err
	}
	res := make([]string, len(body.Steps))
	for i, step := range body.Steps {
		res[i] = step.Name
	}

	return res, nil
}

func GetPipelinesForBranch(client http.PipelineHttpClient, branch string) (*models.PipelinesLookupResponse, error) {
	res, err := requests.GetPipelines(client, models.PipelinesLookupOptions{
		SortBy:   "latestRunId",
		FilterBy: "ja-2446", //should be branch
		Light:    true,
		Limit:    3, // TODO global
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetSyncStatusForBranch(client http.PipelineHttpClient, branch string) (*models.SyncStatusResponse, error) {
	res, err := requests.GetSyncStatus(client, models.SyncOptions{
		PipelineSourceId:       6,
		PipelineSourceBranches: "ja-2446", //should be branch
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetSourcesById(client http.PipelineHttpClient, id string) (*models.SourcesResponse, error) {
	res, err := requests.GetSource(client, models.SourcesOptions{
		PipelineSourceIds: id, // should be 'id'
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetStepConnectionsByPipelinesId(client http.PipelineHttpClient, pipelinesIds string) (*models.StepConnectionsResponse, error) {
	res, err := requests.GetStepConnections(client, models.GetStepConnectionsOptions{
		PipelineIds: pipelinesIds,
		Limit:       3, // TODO const
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetRuns(client http.PipelineHttpClient, pipelinesIds string) (*models.RunsResponse, error) {
	res, err := requests.GetRuns(client, models.GetRunsOptions{
		PipelineIds: pipelinesIds,
		Limit:       3, // TODO const
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func CancelRun(client http.PipelineHttpClient, runId int) (*models.RunsResponse, error) {
	res, err := requests.CancelRun(client, runId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetStepsTestReports(client http.PipelineHttpClient, stepIds string) (*models.GetStepsTestReportResponse, error) {
	res, err := requests.GetStepsTestReports(client, models.StepsTestReportsOptions{
		StepIds: stepIds,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
