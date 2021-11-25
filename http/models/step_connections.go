package models

type GetStepConnectionsOptions struct {
	PipelineIds string `url:"pipelineIds,omitempty"`
	Limit       int    `url:"limit,omitempty"`
}

type StepConnection struct {
	Id                    int    `json:"id"`
	PipelineStepId        int    `json:"pipelineStepId"`
	Operation             string `json:"operation"`
	OperationResourceId   int    `json:"operationResourceId"`
	OperationResourceName string `json:"operationResourceName"`
	PipelineSourceId      int    `json:"pipelineSourceId"`
	PipelineId            int    `json:"pipelineId"`
	IsPipelineConnection  bool   `json:"isPipelineConnection"`
	IsTrigger             bool   `json:"isTrigger"`
}

type StepConnectionsResponse struct {
	StepConnections []StepConnection
}
