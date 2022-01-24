package utils

import (
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"strconv"
	"strings"
)

func GetLatestRunId(httpClient http.PipelineHttpClient, serviceName string, pipeSuffix string, branch string) (latestRunId int) {
	if httpClient == nil || serviceName == "" || pipeSuffix == "" || branch == "" {
		panic("Missing parameters!")
	}
	pipeResp, err := requests.GetPipelines(httpClient, requests.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   branch,
		Light:      true,
		Limit:      1,
		PipesNames: serviceName + pipeSuffix,
	})
	if err != nil {
		panic(err)
	}
	if len(pipeResp.Pipelines) > 0 {
		latestRunId = pipeResp.Pipelines[0].LatestRunId
	}
	return
}

func GetRunDetails(httpClient http.PipelineHttpClient, latestRunId int) (statusCode string, startedAt string, endedAt string, triggeredBy string) {
	if httpClient == nil || latestRunId < 1 {
		panic("Missing parameters!")
	}

	runDetails, err := requests.GetRuns(httpClient, requests.GetRunsOptions{
		RunIds: strconv.Itoa(latestRunId),
	})
	if err != nil {
		panic(err)
	}

	if len(runDetails.Runs) > 1 {
		panic("Unexpected result!")
	}

	if len(runDetails.Runs) == 1 {
		if runDetails.Runs[0].StaticPropertyBag.TriggeredByUserName != "" {
			triggeredBy = runDetails.Runs[0].StaticPropertyBag.TriggeredByUserName
		} else {
			triggeredBy = runDetails.Runs[0].StaticPropertyBag.TriggeredByResourceName
		}
		return runDetails.Runs[0].StatusCode.StatusCodeName(), runDetails.Runs[0].StartedAt, runDetails.Runs[0].EndedAt, triggeredBy
	} else {
		return
	}
}

func GetMilestoneBranchAndVersion(httpClient http.PipelineHttpClient, latestRunId int) (adHocReleaseBranchName string, serviceVersion string) {
	if httpClient == nil || latestRunId < 1 {
		panic("Missing parameters!")
	}
	adHocReleaseBranchName = ""
	serviceVersion = ""

	steps, err := requests.GetSteps(httpClient, requests.GetStepsOptions{
		RunIds: strconv.Itoa(latestRunId),
		Limit:  15,
	})
	if err != nil {
		panic(err)
	}

	var infraReportLinksStepId int
	for _, step := range steps.Steps {
		if step.Name == InfraReportLinksStep {
			infraReportLinksStepId = step.Id
		}
		// In case of manual step re-trigger, pipeline params will be found here
		for _, variable := range step.ConfigPropertyBag.EnvironmentVariables {
			if adHocReleaseBranchName == "" && variable.Key == InfraReportEnvAdHocReleaseBranchName {
				adHocReleaseBranchName = variable.Value
			}
			if serviceVersion == "" && strings.HasSuffix(variable.Key, InfraReportEnvVersionSuffix) {
				serviceVersion = variable.Value
			}
		}
	}
	if infraReportLinksStepId == 0 || (adHocReleaseBranchName != "" && serviceVersion != "") {
		return
	}

	variables, err := requests.GetStepVariables(httpClient, requests.GetStepVariablesOptions{
		StepIds: infraReportLinksStepId,
	})
	if err != nil {
		panic(err)
	}
	if len(variables.Variables) == 0 {
		return
	}
	for _, entity := range variables.Variables[0].Variables.RunVariable {
		if entity.Key == InfraReportEnvAdHocReleaseBranchName {
			adHocReleaseBranchName = entity.Value
		}
		if strings.HasSuffix(entity.Key, InfraReportEnvVersionSuffix) {
			serviceVersion = entity.Value
		}
	}
	return
}
