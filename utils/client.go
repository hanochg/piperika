package utils

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
)

const (
	pipelineUrlPath = "/pipelines/api/v1"
)

type ClientOptions struct {
	Query interface{}
}

type PipelineHttpClient interface {
	SendGet(endpoint string, options ClientOptions) ([]byte, error)
}

func NewPipelineHttp(c *components.Context) (*pipelineHttpClient, error) {
	details, err := plugins.GetServerDetails(c)
	if err != nil {
		return nil, err
	}
	manager, err := utils.CreateServiceManager(details, 1, false)
	if err != nil {
		return nil, err
	}
	config, err := details.CreateArtAuthConfig()
	if err != nil {
		return nil, err
	}

	return &pipelineHttpClient{
		client:  manager.Client(),
		details: config.CreateHttpClientDetails(),
		baseUrl: getPipelineUrlFromBaseUrl(details.Url),
	}, nil
}

func getPipelineUrlFromBaseUrl(url string) string {
	return url + pipelineUrlPath
}

type pipelineHttpClient struct {
	client  *jfroghttpclient.JfrogHttpClient
	details httputils.HttpClientDetails
	baseUrl string
}

func (s *pipelineHttpClient) SendGet(endpoint string, options ClientOptions) ([]byte, error) {
	url, err := getUrlWithQuery(s.baseUrl+endpoint, options)
	if err != nil {
		return nil, err
	}
	res, resBody, _, err := s.client.SendGet(url, true, &s.details)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected response; status code: %d, message: %s", res.StatusCode, resBody)
	}
	return resBody, nil
}

func getUrlWithQuery(baseUrl string, options ClientOptions) (string, error) {
	url := baseUrl
	if options.Query != nil {
		values, err := query.Values(options)
		if err != nil {
			return "", err
		}
		url += values.Encode()
	}
	return url, nil
}
