package report

import (
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/utils"
	"strconv"
	"sync"
	"time"
)

type ServiceReport struct {
	HttpClient             http.PipelineHttpClient
	ServiceName            string
	ProjectName            string
	ProjectID              int
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
	TriggeredAt string
	ProjectId   int
	ProjectName string
}

var (
	projectIdsNames = sync.Map{}
)

func (sr *ServiceReport) fetchReport() {
	if sr.HttpClient == nil || sr.ServiceName == "" || sr.BaseUrl == "" || sr.MilestoneBranch == "" || sr.ProjectName == "" {
		panic("These variables must be initialized: HttpClient, ServiceName, BaseUrl, MilestoneBranch")
	}

	sr.ProjectID = utils.GetProjectIdByName(sr.HttpClient, sr.ProjectName)
	projectIdsNames.Store(sr.ProjectID, sr.ProjectName)
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
	getPipelineReport(sr.HttpClient, sr.ServiceName, utils.InfraReportReleasePipeSuffix, sr.MilestoneBranch, &sr.ReleasePipeline, sr.ProjectID)
	sr.ReleasePipeline.PipeUrl = getPipelineUrl(sr.BaseUrl, sr.ServiceName+utils.InfraReportReleasePipeSuffix, sr.MilestoneBranch, sr.ProjectName)
}

func (sr *ServiceReport) getBuildPipelineReport() {
	getPipelineReport(sr.HttpClient, sr.ServiceName, utils.InfraReportBuildPipeSuffix, sr.AdHocReleaseBranchName, &sr.BuildPipeline, sr.ProjectID)
	sr.BuildPipeline.PipeUrl = getPipelineUrl(sr.BaseUrl, sr.ServiceName+utils.InfraReportBuildPipeSuffix, sr.AdHocReleaseBranchName, sr.ProjectName)
}

func (sr *ServiceReport) getPostReleasePipeline() {
	getPipelineReport(sr.HttpClient, sr.ServiceName, utils.InfraReportPostReleasePipeSuffix, sr.AdHocReleaseBranchName, &sr.PostReleasePipeline, sr.ProjectID)
	sr.PostReleasePipeline.PipeUrl = getPipelineUrl(sr.BaseUrl, sr.ServiceName+utils.InfraReportPostReleasePipeSuffix, sr.AdHocReleaseBranchName, sr.ProjectName)
}

func getPipelineReport(httpClient http.PipelineHttpClient, serviceName string, suffix string, branch string, result *PipelineResult, projectId int) {
	result.Branch = branch
	result.Name = serviceName + suffix
	if branch != "" {
		result.LatestRunId = utils.GetLatestRunId(httpClient, serviceName, suffix, branch, strconv.Itoa(projectId))
	}
	if result.LatestRunId > 0 {
		runDetails := utils.GetRunDetails(httpClient, result.LatestRunId)
		result.Status = runDetails.StatusCodeName
		result.ProjectId = runDetails.ProjectId
		result.TriggeredAt = runDetails.TriggeredAt
		result.TriggeredBy = runDetails.TriggeredBy
		result.StartTime = runDetails.StartedAt
		result.EndTime = runDetails.EndedAt
	}
}

func getPipelineUrl(uiUrl string, service string, branchName string, project string) string {
	if branchName == "" {
		return "Pipeline is not created."
	} else {
		return fmt.Sprintf("%s ", utils.GetPipelinesBranchURL(uiUrl, service, branchName, project))
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

func (pr *PipelineResult) string() string {
	return fmt.Sprintf(
		"%s %s%s: Status: %s, TriggeredBy: %s, TriggeredAt: %s, StartTime: %s, EndTime: %s"+
			"\n%s",
		pr.statusIcon(), pr.projectNamePrefix(), pr.Name, pr.Status, pr.TriggeredBy, pr.TriggeredAt, pr.StartTime, pr.EndTime, pr.PipeUrl)
}

func (pr *PipelineResult) statusIcon() string {
	triggerTime, err := time.Parse(time.RFC3339, pr.TriggeredAt)
	tonightTime := time.Now().Truncate(24 * time.Hour)

	if pr.Status == "Success" && err == nil {
		if triggerTime.After(tonightTime) {
			return "✅"
		} else {
			return "⏰"
		}
	} else {
		return "❌"
	}
}

func (pr *PipelineResult) projectNamePrefix() string {
	actualProjectName, ok := projectIdsNames.Load(pr.ProjectId)
	if ok {
		return fmt.Sprint(actualProjectName) + "/"
	} else {
		return ""
	}
}
