package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	resourcesUrl = "/resourceVersions"
)

type GetResourcesOptions struct {
	PipelineSourceIds  string `url:"pipelineSourceIds,omitempty"`  // Can be a csv list
	ResourceVersionIds string `url:"resourceVersionIds,omitempty"` // Can be a csv list
	RunIds             string `url:"runIds,omitempty"`             // Can be a csv list
}

type ContentPropertyBag struct {
	Path       string `json:"path"`
	CommitSha  string `json:"commitSha"`
	BranchName string `json:"branchName"`
}

type Resource struct {
	ContentPropertyBag ContentPropertyBag `json:"contentPropertyBag"`
	Id                 int                `json:"id"`
	ResourceTypeCode   http.ResourceCodes `json:"resourceTypeCode"`
}

type ResourcesResponse struct {
	Resources []Resource
}

func GetResourceVersions(client http.PipelineHttpClient, options GetResourcesOptions) (*ResourcesResponse, error) {
	body, err := client.SendGet(resourcesUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &ResourcesResponse{}
	err = json.Unmarshal(body, &res.Resources)
	return res, err
}
