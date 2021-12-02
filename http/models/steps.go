package models

type GetStepsOptions struct {
	RunIds string `url:"runIds,omitempty"`
	Limit  int    `url:"limit,omitempty"`
}

type Step struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	StatusCode StatusCode `json:"statusCode"`
}

type StepsResponse struct {
	Steps []Step
}
