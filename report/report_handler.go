package report

import (
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"strconv"
	"strings"
)

const adHocReleaseBranchName = "AD_HOC_RELEASE_BRANCH_NAME"

func getPipelineReport(httpClient http.PipelineHttpClient, serviceName string, pipeSuffix string, branch string, pipeline *PipelineResult) {
	pipeline.Name = serviceName + pipeSuffix
	pipeline.Branch = branch

	pipeResp, _ := requests.GetPipelines(httpClient, requests.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   branch,
		Light:      true,
		Limit:      1,
		PipesNames: pipeline.Name,
	})
	latestRunId := pipeResp.Pipelines[0].LatestRunId
	if latestRunId == 0 {
		pipeline.Status = "NotTriggered"
		return
	}

	pipeline.LatestRunId = latestRunId
	runDetails, _ := requests.GetRuns(httpClient, requests.GetRunsOptions{
		RunIds: strconv.Itoa(latestRunId),
	})
	pipeline.Status = runDetails.Runs[0].StatusCode.StatusCodeName()
	pipeline.StartTime = runDetails.Runs[0].StartedAt
	pipeline.EndTime = runDetails.Runs[0].EndedAt
	triggeredByUser := runDetails.Runs[0].StaticPropertyBag.TriggeredByUserName
	if triggeredByUser != "" {
		pipeline.TriggeredBy = triggeredByUser
	} else {
		pipeline.TriggeredBy = runDetails.Runs[0].StaticPropertyBag.TriggeredByResourceName
	}
}

func getMilestoneBranchAndVersion(httpClient http.PipelineHttpClient, latestRunId int) (string, string) {
	steps, _ := requests.GetSteps(httpClient, requests.GetStepsOptions{
		RunIds: strconv.Itoa(latestRunId),
	})
	releaseProcessLinksStepId := -1
	for _, step := range steps.Steps {
		if step.Name == "release_process_links" {
			releaseProcessLinksStepId = step.Id
		}
	}

	variables, _ := requests.GetStepVariables(httpClient, requests.GetStepVariablesOptions{
		StepIds: releaseProcessLinksStepId,
	})

	adHocReleaseBranchNameValue := ""
	serviceVersion := ""
	for _, entity := range variables.Variables[0].Variables.RunVariable {
		if entity.Key == adHocReleaseBranchName {
			adHocReleaseBranchNameValue = entity.Value
		}
		if strings.HasSuffix(entity.Key, "_VERSION") {
			serviceVersion = entity.Value
		}
	}
	return adHocReleaseBranchNameValue, serviceVersion
}
