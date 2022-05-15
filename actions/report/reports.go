package report

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ServiceReport struct {
	httpClient             http.PipelineHttpClient
	serviceName            string
	BaseUrl                string
	milestoneBranch        string
	releasePipeline        PipelineResult
	buildPipeline          PipelineResult
	postReleasePipeline    PipelineResult
	adHocReleaseBranchName string
	serviceVersion         string
	projectName            string
	config                 *utils.Reports
}

type PipelineResult struct {
	name        string
	branch      string
	pipeUrl     string
	latestRunId int
	StartTime   string
	endTime     string
	status      string
	triggeredBy string
}

func ReportsGathering(ctx context.Context) error {
	config := ctx.Value(utils.ConfigCtxKey).(*utils.Configurations)
	branchName := ctx.Value(utils.BranchName).(string)
	baseUiUrl := ctx.Value(utils.BaseUiUrl).(string)
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	err := validateConfigurations(config)
	if err != nil {
		return err
	}

	reportsToPrint := buildServicesReportStructs(httpClient, baseUiUrl, branchName, config)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for i := 0; i < len(reportsToPrint); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			curReport := reportsToPrint[i]
			err := curReport.fetchReport()
			if err != nil {
				fmt.Println(fmt.Sprintf("could not fetch report for service %s, err %v", curReport.serviceName, err))
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			fmt.Println(curReport.string())
		}(i)
	}
	wg.Wait()
	fmt.Println("Finished")
	return nil
}

func validateConfigurations(config *utils.Configurations) error {
	if config.Reports == nil || len(config.Reports.ServicesNameAndProject) == 0 ||
		config.Reports.VersionSuffix == "" || config.Reports.BuildPipeSuffix == "" ||
		config.Reports.PostReleasePipeSuffix == "" || config.Reports.ReleasePipeSuffix == "" ||
		config.Reports.AdHocReleaseBranchName == "" || config.Reports.AdHocReleaseBranchLinksStep == "" {
		return fmt.Errorf("invalid report configurations were provided, exiting")
	}
	return nil
}

func buildServicesReportStructs(httpClient http.PipelineHttpClient, baseUrl string, branch string, config *utils.Configurations) []*ServiceReport {
	var rep []*ServiceReport
	for curPipeName, curPipeProjectName := range config.Reports.ServicesNameAndProject {
		el := &ServiceReport{
			httpClient:      httpClient,
			BaseUrl:         baseUrl,
			config:          config.Reports,
			milestoneBranch: branch,
			serviceName:     curPipeName,
			projectName:     curPipeProjectName,
		}
		rep = append(rep, el)
	}
	return rep
}

func (sr *ServiceReport) fetchReport() error {
	if sr.serviceName == "" || sr.BaseUrl == "" || sr.milestoneBranch == "" {
		return fmt.Errorf("these variables must be initialized: ServiceName, BaseUrl, MilestoneBranch")
	}

	err := sr.getReleasePipelineReport()
	if err != nil {
		return err
	}
	if sr.releasePipeline.latestRunId > 0 {
		var err error
		sr.adHocReleaseBranchName, sr.serviceVersion, err = getMilestoneBranchAndVersion(sr.httpClient, sr.config, sr.releasePipeline.latestRunId)
		if err != nil {
			return err
		}
	}

	err = sr.getBuildPipelineReport()
	if err != nil {
		return err
	}
	err = sr.getPostReleasePipeline()
	if err != nil {
		return err
	}
	return nil
}

func (sr *ServiceReport) getReleasePipelineReport() error {
	err := getPipelineReport(sr.httpClient, sr.serviceName, sr.config.ReleasePipeSuffix, sr.milestoneBranch, &sr.releasePipeline)
	if err != nil {
		return err
	}
	sr.releasePipeline.pipeUrl = utils.GetPipelinesBranchURL(sr.BaseUrl, sr.serviceName, sr.config.ReleasePipeSuffix, sr.milestoneBranch, sr.projectName)
	return nil
}

func (sr *ServiceReport) getBuildPipelineReport() error {
	err := getPipelineReport(sr.httpClient, sr.serviceName, sr.config.BuildPipeSuffix, sr.adHocReleaseBranchName, &sr.buildPipeline)
	if err != nil {
		return err
	}
	sr.buildPipeline.pipeUrl = utils.GetPipelinesBranchURL(sr.BaseUrl, sr.serviceName, sr.config.BuildPipeSuffix, sr.adHocReleaseBranchName, sr.projectName)
	return nil
}

func (sr *ServiceReport) getPostReleasePipeline() error {
	err := getPipelineReport(sr.httpClient, sr.serviceName, sr.config.PostReleasePipeSuffix, sr.adHocReleaseBranchName, &sr.postReleasePipeline)
	if err != nil {
		return err
	}
	sr.postReleasePipeline.pipeUrl = utils.GetPipelinesBranchURL(sr.BaseUrl, sr.serviceName, sr.config.PostReleasePipeSuffix, sr.adHocReleaseBranchName, sr.projectName)
	return nil
}

func getPipelineReport(httpClient http.PipelineHttpClient, serviceName string, suffix string, branch string, result *PipelineResult) error {
	result.branch = branch
	result.name = serviceName + suffix
	var err error
	if branch != "" {
		result.latestRunId, err = getLatestRunId(httpClient, serviceName+suffix, branch)
		if err != nil {
			return err
		}
	}
	if result.latestRunId > 0 {
		result.status, result.StartTime, result.endTime, result.triggeredBy, err = getRunDetails(httpClient, result.latestRunId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sr ServiceReport) string() string {
	return fmt.Sprintf(
		"\n\t****   %s: %s   ****\n"+
			"%s\n"+
			"%s\n"+
			"%s\n",
		sr.serviceName, sr.serviceVersion, sr.releasePipeline.string(), sr.buildPipeline.string(), sr.postReleasePipeline.string())
}

func (pr PipelineResult) string() string {
	return fmt.Sprintf(
		"%s %s: Status: %s, TriggeredBy: %s, StartTime: %s, EndTime: %s"+
			"\n%s",
		statusIcon(pr.status, pr.StartTime), pr.name, pr.status, pr.triggeredBy, pr.StartTime, pr.endTime, pr.pipeUrl)
}

func statusIcon(status string, startTime string) string {
	triggerTime, err := time.Parse(time.RFC3339, startTime)
	tonightTime := time.Now().Truncate(24 * time.Hour)

	if status == "Success" && err == nil {
		if triggerTime.After(tonightTime) {
			return "✅"
		} else {
			return "⏰"
		}
	} else {
		return "❌"
	}
}

func getLatestRunId(httpClient http.PipelineHttpClient, pipeName string, branch string) (latestRunId int, err error) {
	if httpClient == nil || pipeName == "" || branch == "" {
		panic("Missing parameters!")
	}
	pipeResp, err := requests.GetPipelines(httpClient, requests.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   branch,
		Light:      true,
		Limit:      1,
		PipesNames: pipeName,
	})
	if err != nil {
		return
	}
	if len(pipeResp.Pipelines) > 0 {
		latestRunId = pipeResp.Pipelines[0].LatestRunId
	}
	return
}

func getRunDetails(httpClient http.PipelineHttpClient, latestRunId int) (statusCode string, startedAt string, endedAt string, triggeredBy string, err error) {
	if httpClient == nil || latestRunId < 1 {
		err = fmt.Errorf("missing parameters, could not fetch run details")
		return
	}

	runDetails, err := requests.GetRuns(httpClient, requests.GetRunsOptions{
		RunIds: strconv.Itoa(latestRunId),
	})
	if err != nil {
		return
	}

	if len(runDetails.Runs) > 1 {
		err = fmt.Errorf("unexpected result, could not fetch run details")
		return
	}

	if len(runDetails.Runs) == 1 {
		if runDetails.Runs[0].StaticPropertyBag.TriggeredByUserName != "" {
			triggeredBy = runDetails.Runs[0].StaticPropertyBag.TriggeredByUserName
		} else {
			triggeredBy = runDetails.Runs[0].StaticPropertyBag.TriggeredByResourceName
		}
		statusCode = runDetails.Runs[0].StatusCode.StatusCodeName()
		startedAt = runDetails.Runs[0].StartedAt
		endedAt = runDetails.Runs[0].EndedAt
	}

	return
}

func getMilestoneBranchAndVersion(httpClient http.PipelineHttpClient, config *utils.Reports, latestRunId int) (adHocReleaseBranchName string, serviceVersion string, err error) {
	if httpClient == nil || latestRunId < 1 {
		err = fmt.Errorf("missing parameters, could not get milestone branch and version")
		return
	}
	adHocReleaseBranchName = ""
	serviceVersion = ""

	steps, err := requests.GetSteps(httpClient, requests.GetStepsOptions{
		RunIds: strconv.Itoa(latestRunId),
		Limit:  15,
	})
	if err != nil {
		return
	}

	var reportLinksStepId int
	for _, step := range steps.Steps {
		if step.Name == config.AdHocReleaseBranchLinksStep {
			reportLinksStepId = step.Id
		}
		// In case of manual step re-trigger, pipeline params will be found here
		for _, variable := range step.ConfigPropertyBag.EnvironmentVariables {
			if adHocReleaseBranchName == "" && variable.Key == config.AdHocReleaseBranchName {
				adHocReleaseBranchName = variable.Value
			}
			if serviceVersion == "" && strings.HasSuffix(variable.Key, config.VersionSuffix) {
				serviceVersion = variable.Value
			}
		}
	}
	if reportLinksStepId == 0 || (adHocReleaseBranchName != "" && serviceVersion != "") {
		return
	}

	variables, err := requests.GetStepVariables(httpClient, requests.GetStepVariablesOptions{
		StepIds: reportLinksStepId,
	})
	if err != nil {
		return
	}
	if len(variables.Variables) == 0 {
		return
	}
	for _, entity := range variables.Variables[0].Variables.RunVariable {
		if entity.Key == config.AdHocReleaseBranchName {
			adHocReleaseBranchName = entity.Value
		}
		if strings.HasSuffix(entity.Key, config.VersionSuffix) {
			serviceVersion = entity.Value
		}
	}
	return
}
