package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

//asdsaasd
const (
	runResourcesUrl = "/runResourceVersions"
)

type GetRunResourcesOptions struct {
	Limit             int    `url:"limit,omitempty"`             // Can be a csv list
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"` // Can be a csv list
	RunIds            string `url:"runIds,omitempty"`            // Can be a csv list
	SortBy            string `url:"sortBy,omitempty"`
	SortOrder         int    `url:"sortOrder,omitempty"`
}

type RunContentPropertyBag struct {
	Path       string `json:"path"`
	CommitSha  string `json:"commitSha"`
	BranchName string `json:"branchName"`
}

type RunResource struct {
	ResourceVersionContentPropertyBag RunContentPropertyBag `json:"resourceVersionContentPropertyBag"`
	Id                                int                   `json:"id"`
	RunId                             int                   `json:"runId"`
	ResourceName                      string                `json:"resourceName"`
	ResourceTypeCode                  http.ResourceCodes    `json:"resourceTypeCode"`
}

type RunResourcesResponse struct {
	Resources []RunResource
}

func GetRunResourceVersions(client http.PipelineHttpClient, options GetRunResourcesOptions) (*RunResourcesResponse, error) {
	body, err := client.SendGet(runResourcesUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &RunResourcesResponse{}
	err = json.Unmarshal(body, &res.Resources)
	return res, err
}
