package utils

import (
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strconv"
	"strings"
)

type RunDetails struct {
	StatusCodeName string
	ProjectId      int
	TriggeredAt    string
	TriggeredBy    string
	StartedAt      string
	EndedAt        string
}

func GetProjectIdByName(httpClient http.PipelineHttpClient, projectName string) (projectId int) {
	if projectName == DefaultProject {
		return DefaultProjectId
	}
	projectsResponse, err := requests.GetProjects(httpClient, requests.GetProjectsOptions{
		Names: projectName,
	})
	if err != nil {
		panic(err)
	}
	if len(projectsResponse.Projects) == 1 {
		return projectsResponse.Projects[0].ProjectIds
	} else {
		panic(fmt.Sprintf("Project '%s' not found.", projectName))
	}
}

func GetLatestRunId(httpClient http.PipelineHttpClient, serviceName string, pipeSuffix string, branch string, projectIds string) (latestRunId int) {
	if httpClient == nil || serviceName == "" || pipeSuffix == "" || branch == "" || projectIds == "" {
		panic("Missing parameters!")
	}
	pipeResp, err := requests.GetPipelines(httpClient, requests.GetPipelinesOptions{
		SortBy:                 "latestRunId",
		PipelineSourceBranches: branch,
		ProjectIds:             projectIds,
		Light:                  true,
		Limit:                  5,
		SortOrder:              -1,
		PipesNames:             serviceName + pipeSuffix,
	})
	if err != nil {
		panic(err)
	}
	responseLength := len(pipeResp.Pipelines)
	if responseLength < 2 {
		latestRunId = pipeResp.Pipelines[0].LatestRunId
	} else {
		log.Error(fmt.Sprintf("GetLatestRunId must return exactly one pipelineId, but you've got at least %d pipelines. "+
			"Check you query!!!\nQuery response: %v", responseLength, pipeResp))
	}
	return
}

func GetRunDetails(httpClient http.PipelineHttpClient, latestRunId int) (run RunDetails) {
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

	run = RunDetails{}
	if len(runDetails.Runs) == 1 {
		if runDetails.Runs[0].StaticPropertyBag.TriggeredByUserName != "" {
			run.TriggeredBy = runDetails.Runs[0].StaticPropertyBag.TriggeredByUserName
		} else {
			run.TriggeredBy = runDetails.Runs[0].StaticPropertyBag.TriggeredByResourceName
		}
		run.StatusCodeName = runDetails.Runs[0].StatusCode.StatusCodeName()
		run.ProjectId = runDetails.Runs[0].ProjectId
		run.TriggeredAt = runDetails.Runs[0].TriggeredAt
		run.StartedAt = runDetails.Runs[0].StartedAt
		run.EndedAt = runDetails.Runs[0].EndedAt
	}
	return
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
