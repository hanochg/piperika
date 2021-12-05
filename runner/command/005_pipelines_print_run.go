package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
	"strings"
)

func New005PipelinesPrintRun() *_005 {
	return &_005{}
}

type _005 struct{}

func (c *_005) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	runCompleted := false

	runStatus, err := requests.GetRuns(httpClient, models.GetRunsOptions{
		RunIds: strconv.Itoa(state.RunId),
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline runs data: %v", err),
		}
	}
	if len(runStatus.Runs) == 0 {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline runs data: %v", err),
		}
	}

	statusCode := runStatus.Runs[0].StatusCode
	if statusCode != models.Creating && statusCode != models.Waiting && statusCode != models.Processing {
		// Run completed
		runCompleted = true
	}

	// Get steps statuses
	steps, err := requests.GetSteps(httpClient, models.GetStepsOptions{
		RunIds: strconv.Itoa(state.RunId),
		Limit:  0,
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps data for run id '%d': %v", state.RunId, err),
		}
	}

	stepsIdToNames := make(map[int]string, 0)
	failedSteps := make([]string, 0)
	successSteps := make([]string, 0)
	processingSteps := make([]string, 0)
	for _, step := range steps.Steps {
		stepsIdToNames[step.Id] = step.Name
		if step.StatusCode == models.Failure {
			failedSteps = append(failedSteps, step.Name)
		}
		if step.StatusCode == models.Success {
			successSteps = append(successSteps, step.Name)
		}
		if step.StatusCode == models.Processing {
			processingSteps = append(processingSteps, step.Name)
		}
	}

	if !runCompleted {
		outputMsg := fmt.Sprintf("Run number %d - Completed %d out of %d. Steps (InProgress/Succeed/Failed/Total) %d/%d/%d/%d",
			state.RunNumber, len(failedSteps)+len(successSteps), len(steps.Steps), len(processingSteps),
			len(successSteps), len(failedSteps), len(steps.Steps))
		if len(failedSteps) != 0 {
			outputMsg += fmt.Sprintf("\nFailed steps: %s", strings.Join(failedSteps, ","))
		}

		return Status{
			PipelinesStatus: "processing",
			Message:         outputMsg,
			Type:            InProgress,
		}
	}

	testsFailureOutput, err := c.createTestReport(httpClient, state, stepsIdToNames)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps test reports for run id '%d': %v", state.RunId, err),
		}
	}

	outputStr := fmt.Sprintf("Run %d completed - (InProgress/Succeed/Failed) %d/%d/%d.\nFailed steps: %s\nTests results: %s",
		state.RunNumber, len(steps.Steps), len(failedSteps), len(successSteps), strings.Join(failedSteps, ","), testsFailureOutput)
	if testsFailureOutput != "" {
		outputStr += fmt.Sprintf("\nTests results: %s", testsFailureOutput)
	}
	return Status{
		Message: outputStr,
		Type:    Done,
	}
}

func (c *_005) createTestReport(httpClient http.PipelineHttpClient, state *PipedCommandState, stepsIdToNames map[int]string) (string, error) {
	testReports, err := requests.GetStepsTestReports(httpClient, models.StepsTestReportsOptions{StepIds: state.RunStepIdsCsv})
	if err != nil {
		return "", err
	}

	var testsFailureOutput bytes.Buffer
	for _, testFailures := range testReports.TestReports {
		for _, failure := range testFailures.FailureDetails {
			testsFailureOutput.WriteString(fmt.Sprintf("\n%s - %s:%s - %s, %s",
				failure.ClassName, failure.TestName, failure.Kind, failure.Message, stepsIdToNames[testFailures.StepId]))
		}
	}
	return testsFailureOutput.String(), nil
}

func (c *_005) TriggerOnFail(_ context.Context, _ *PipedCommandState) error {
	return fmt.Errorf("timed out")
}
