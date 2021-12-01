package models

type GetResourcesOptions struct {
	PipelineSourceIds  string `url:"pipelineSourceIds,omitempty"`  // Can be a csv list
	ResourceVersionIds string `url:"resourceVersionIds,omitempty"` // Can be a csv list
}

type ContentPropertyBag struct {
	Path       string `json:"path"`
	CommitSha  string `json:"commitSha"`
	BranchName string `json:"branchName"`
}

type Resource struct {
	ContentPropertyBag ContentPropertyBag `json:"contentPropertyBag"`
	Id                 int                `json:"id"`
}

type ResourcesResponse struct {
	Resources []Resource
}
