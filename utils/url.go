package utils

import (
	"fmt"
	"github.com/buger/goterm"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"net/url"
)

const (
	baseUiUrl = "/ui/pipelines"
)

func GetUIBaseUrl(c *components.Context) (string, error) {
	details, err := plugins.GetServerDetails(c)
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

func GetPipelinesRunURL(uiBaseUrl string, pipelineName string, step string, runNumber int, gitBranch string) string {
	return goterm.Color(fmt.Sprintf("%s/myPipelines/default/%s/%d/%s?branch=%v", uiBaseUrl, pipelineName, runNumber, step, url.PathEscape(gitBranch)),
		goterm.BLUE)
}

func GetPipelinesBranchURL(uiBaseUrl string, pipelineName string, suffix string, gitBranch string) string {
	pipelinesNameWithSuffix := pipelineName + suffix
	if gitBranch == "" {
		return goterm.Color(fmt.Sprintf("Pipeline '%s' was not created", pipelinesNameWithSuffix),
			goterm.RED)
	}
	return goterm.Color(
		fmt.Sprintf("%s/myPipelines/default/%s?branch=%v",
			uiBaseUrl, pipelinesNameWithSuffix, url.PathEscape(gitBranch)),
		goterm.BLUE)
}
