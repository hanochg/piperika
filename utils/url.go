package utils

import (
	"fmt"
	"github.com/buger/goterm"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/common"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"net/url"
)

const (
	baseUiUrl = "/ui/pipelines"
)

func GetProjectNameForSource(httpClient http.PipelineHttpClient, pipelinesSourceId int) (string, error) {
	resp, err := requests.SyncOrGetSource(httpClient, requests.SyncSourcesOptions{
		PipelineSourceId: pipelinesSourceId,
	})
	if err != nil {
		return "", nil
	}
	if len(resp.Sources) == 0 {
		return "", fmt.Errorf("could not get project name for source %d", pipelinesSourceId)
	}

	curProjectId := resp.Sources[0].ProjectId
	projResp, err := requests.GetProjects(httpClient, requests.ProjectsOptions{
		ProjectId: curProjectId,
	})
	if err != nil {
		return "", err
	}
	return projResp.Name, nil
}

func GetUIBaseUrl(c *components.Context) (string, error) {
	details, err := common.GetServerDetails(c)
	if err != nil {
		return "", err
	}

	baseUrl, err := url.Parse(details.Url)
	if err != nil {
		return "", err
	}

	basePath, err := url.Parse(baseUiUrl)
	if err != nil {
		return "", err
	}

	return baseUrl.ResolveReference(basePath).String(), nil
}

func GetPipelinesRunURL(uiBaseUrl string, pipelineName string, step string, runNumber int, gitBranch string, projectName string) string {
	return goterm.Color(fmt.Sprintf("%s/myPipelines/%s/%s/%d/%s?branch=%v", uiBaseUrl, projectName, pipelineName, runNumber, step, url.PathEscape(gitBranch)),
		goterm.BLUE)
}

func GetPipelinesBranchURL(uiBaseUrl string, pipelineName string, suffix string, gitBranch string, projectName string) string {
	pipelinesNameWithSuffix := pipelineName + suffix
	if gitBranch == "" {
		return goterm.Color(fmt.Sprintf("Pipeline '%s' was not created", pipelinesNameWithSuffix),
			goterm.RED)
	}
	return goterm.Color(
		fmt.Sprintf("%s/myPipelines/%s/%s?branch=%v",
			uiBaseUrl, projectName, pipelinesNameWithSuffix, url.PathEscape(gitBranch)),
		goterm.BLUE)
}
