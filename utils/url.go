package utils

import (
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
