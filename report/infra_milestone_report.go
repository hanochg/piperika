package report

import (
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/utils"
	"time"
)

type ServiceReport struct {
	HttpClient             http.PipelineHttpClient
	ServiceName            string
	BaseUrl                string
	MilestoneBranch        string
	ReleasePipeline        PipelineResult
	BuildPipeline          PipelineResult
	PostReleasePipeline    PipelineResult
	AdHocReleaseBranchName string
	ServiceVersion         string
}

type PipelineResult struct {
	Name        string
	Branch      string
	PipeUrl     string
	LatestRunId int
	StartTime   string
	EndTime     string
	Status      string
	TriggeredBy string
}

func (sr *ServiceReport) fetchReport() {
	if sr.HttpClient == nil || sr.ServiceName == "" || sr.BaseUrl == "" || sr.MilestoneBranch == "" {
		panic("These variables must be initialized: HttpClient, ServiceName, BaseUrl, MilestoneBranch")
	}
	sr.getReleasePipelineReport()
	if sr.ReleasePipeline.LatestRunId > 0 {
		sr.AdHocReleaseBranchName, sr.ServiceVersion = utils.GetMilestoneBranchAndVersion(sr.HttpClient, sr.ReleasePipeline.LatestRunId)
	}
	sr.getBuildPipelineReport()
	sr.getPostReleasePipeline()
}

func (sr *ServiceReport) toString() string {
	return sr.string()
}

func (sr *ServiceReport) getReleasePipelineReport() {
	getPipelineReport(sr.HttpClient, sr.ServiceName, utils.InfraReportReleasePipeSuffix, sr.MilestoneBranch, &sr.ReleasePipeline)
	sr.ReleasePipeline.PipeUrl = getPipelineUrl(sr.BaseUrl, sr.ServiceName+utils.InfraReportReleasePipeSuffix, sr.MilestoneBranch)
}

func (sr *ServiceReport) getBuildPipelineReport() {
	getPipelineReport(sr.HttpClient, sr.ServiceName, utils.InfraReportBuildPipeSuffix, sr.AdHocReleaseBranchName, &sr.BuildPipeline)
	sr.BuildPipeline.PipeUrl = getPipelineUrl(sr.BaseUrl, sr.ServiceName+utils.InfraReportBuildPipeSuffix, sr.AdHocReleaseBranchName)
}

func (sr *ServiceReport) getPostReleasePipeline() {
	getPipelineReport(sr.HttpClient, sr.ServiceName, utils.InfraReportPostReleasePipeSuffix, sr.AdHocReleaseBranchName, &sr.PostReleasePipeline)
	sr.PostReleasePipeline.PipeUrl = getPipelineUrl(sr.BaseUrl, sr.ServiceName+utils.InfraReportPostReleasePipeSuffix, sr.AdHocReleaseBranchName)
}

func getPipelineReport(httpClient http.PipelineHttpClient, serviceName string, suffix string, branch string, result *PipelineResult) {
	result.Branch = branch
	result.Name = serviceName + suffix
	if branch != "" {
		result.LatestRunId = utils.GetLatestRunId(httpClient, serviceName, suffix, branch)
	}
	if result.LatestRunId > 0 {
		result.Status, result.StartTime, result.EndTime, result.TriggeredBy = utils.GetRunDetails(httpClient, result.LatestRunId)
	}
}

func getPipelineUrl(uiUrl string, service string, branchName string) string {
	if branchName == "" {
		return "Pipeline is not created."
	} else {
		return fmt.Sprintf("%s ", utils.GetPipelinesBranchURL(uiUrl, service, branchName))
	}
}

func (sr ServiceReport) string() string {
	return fmt.Sprintf(
		"\n\t****   %s: %s   ****\n"+
			"%s\n"+
			"%s\n"+
			"%s\n",
		sr.ServiceName, sr.ServiceVersion, sr.ReleasePipeline.string(), sr.BuildPipeline.string(), sr.PostReleasePipeline.string())
}

func (pr PipelineResult) string() string {
	return fmt.Sprintf(
		"%s %s: Status: %s, TriggeredBy: %s, StartTime: %s, EndTime: %s"+
			"\n%s",
		statusIcon(pr.Status, pr.StartTime), pr.Name, pr.Status, pr.TriggeredBy, pr.StartTime, pr.EndTime, pr.PipeUrl)
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
