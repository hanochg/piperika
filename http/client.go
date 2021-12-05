package http

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"net/url"
	"path"
)

const (
	apiV1 = "/api/v1"
)

type ClientOptions struct {
	Query interface{}
}

type PipelineHttpClient interface {
	SendGet(endpoint string, options ClientOptions) ([]byte, error)
	SendPost(endpoint string, options ClientOptions, content []byte) ([]byte, error)
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
	baseUrl, err := getBaseUrlPath(details)
	if err != nil {
		return nil, err
	}
	return &pipelineHttpClient{
		client:  manager.Client(),
		details: config.CreateHttpClientDetails(),
		baseUrl: baseUrl,
	}, nil
}

func getBaseUrlPath(details *config.ServerDetails) (string, error) {
	baseUrl, err := url.Parse(details.PipelinesUrl)
	if err != nil {
		return "", err
	}

	baseUrl.Path = path.Join(baseUrl.Path, apiV1)
	return baseUrl.String(), nil
}

type pipelineHttpClient struct {
	client  *jfroghttpclient.JfrogHttpClient
	details httputils.HttpClientDetails
	baseUrl string
}

func (s *pipelineHttpClient) SendPost(endpoint string, options ClientOptions, content []byte) ([]byte, error) {
	url, err := getUrlWithQuery(s.baseUrl+endpoint, options)
	if err != nil {
		return nil, err
	}
	res, resBody, err := s.client.SendPost(url, content, &s.details)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected response; status code: %d, from: %s, message: %s", res.StatusCode, url, resBody)
	}
	return resBody, nil
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
		return nil, fmt.Errorf("unexpected response; status code: %d, from: %s, message: %s", res.StatusCode, url, resBody)
	}
	return resBody, nil
}

func getUrlWithQuery(baseUrl string, options ClientOptions) (string, error) {
	url := baseUrl
	if options.Query != nil {
		values, err := query.Values(options.Query)
		if err != nil {
			return "", err
		}
		url += "?" + values.Encode()
	}
	return url, nil
}
