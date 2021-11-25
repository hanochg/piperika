package models

type StepsTestReportsOptions struct {
	StepIds string `url:"stepIds,omitempty"`
}

type TestDetails struct {
	TestName    string `json:"testName"`
	ClassName   string `json:"className"`
	SuiteName   string `json:"suiteName"`
	Message     string `json:"message"`
	Kind        string `json:"kind"`
	FullMessage string `json:"full"`
}

type StepTestReport struct {
	Id               int           `json:"id"`
	PipelineSourceId int           `json:"pipelineSourceId"`
	StepId           int           `json:"stepId"`
	DurationSeconds  int           `json:"durationSeconds"`
	TotalTests       int           `json:"totalTests"`
	TotalPassing     int           `json:"totalPassing"`
	TotalFailures    int           `json:"totalFailures"`
	TotalErrors      int           `json:"totalErrors"`
	TotalSkipped     int           `json:"totalSkipped"`
	ErrorDetails     []TestDetails `json:"errorDetails"`
	FailureDetails   []TestDetails `json:"failureDetails"`
}

type GetStepsTestReportResponse struct {
	TestReports []StepTestReport
}
