package report

import (
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/utils"
	"time"
)

type ServiceReport struct {
	ServiceName            string
	ReleasePipeline        PipelineResult
	BuildPipeline          PipelineResult
	PostReleasePipeline    PipelineResult
	ReleaseBranchName      string
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

func GetServiceReport(client http.PipelineHttpClient, serviceName string, uiUrl string, branch string) ServiceReport {
	serviceReport := ServiceReport{}
	serviceReport.ServiceName = serviceName
	serviceReport.ReleaseBranchName = branch

	getReleasePipelineReport(uiUrl, client, serviceName, branch, &serviceReport.ReleasePipeline)
	serviceReport.AdHocReleaseBranchName, serviceReport.ServiceVersion = getMilestoneBranchAndVersion(client, serviceReport.ReleasePipeline.LatestRunId)
	getBuildPipelineReport(uiUrl, client, serviceName, serviceReport.AdHocReleaseBranchName, &serviceReport.BuildPipeline)
	getPostReleasePipeline(uiUrl, client, serviceName, serviceReport.AdHocReleaseBranchName, &serviceReport.PostReleasePipeline)

	return serviceReport
}

func getReleasePipelineReport(uiUrl string, httpClient http.PipelineHttpClient, serviceName string, branch string, result *PipelineResult) {
	const suffix = "_release"
	getPipelineReport(httpClient, serviceName, suffix, branch, result)
	result.PipeUrl = getPipelineUrl(uiUrl, serviceName+suffix, branch)
}

func getBuildPipelineReport(uiUrl string, httpClient http.PipelineHttpClient, serviceName string, branch string, result *PipelineResult) {
	const suffix = "_build"
	getPipelineReport(httpClient, serviceName, suffix, branch, result)
	result.PipeUrl = getPipelineUrl(uiUrl, serviceName+suffix, branch)
}

func getPostReleasePipeline(uiUrl string, httpClient http.PipelineHttpClient, serviceName string, branch string, result *PipelineResult) {
	const suffix = "_post_release"
	getPipelineReport(httpClient, serviceName, suffix, branch, result)
	result.PipeUrl = getPipelineUrl(uiUrl, serviceName+suffix, branch)
}

func getPipelineUrl(uiUrl string, service string, branchName string) string {
	return fmt.Sprintf("%s ",
		utils.GetPipelinesBranchURL(uiUrl, service, branchName))
}

func (r ServiceReport) String() string {
	return fmt.Sprintf(
		"\n\t****   %s: %s   ****\n"+
			"%s\n"+
			"%s\n"+
			"%s\n",
		r.ServiceName, r.ServiceVersion, r.ReleasePipeline.String(), r.BuildPipeline.String(), r.PostReleasePipeline.String())
}

func (p PipelineResult) String() string {
	return fmt.Sprintf(
		"%s %s: Status: %s, TriggeredBy: %s, StartTime: %s, EndTime: %s"+
			"\n%s",
		statusIcon(p.Status, p.StartTime), p.Name, p.Status, p.TriggeredBy, p.StartTime, p.EndTime, p.PipeUrl)
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
