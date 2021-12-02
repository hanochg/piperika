package models

type GetRunResourcesOptions struct {
	PipelineSourceIds string `url:"pipelineSourceIds,omitempty"` // Can be a csv list
	RunIds            string `url:"runIds,omitempty"`            // Can be a csv list
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
	ResourceTypeCode                  ResourceTypes         `json:"resourceTypeCode"`
}

type RunResourcesResponse struct {
	Resources []RunResource
}
