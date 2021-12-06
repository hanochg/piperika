package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/buger/goterm"
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
	dirConfig := ctx.Value(utils.DirConfigCtxKey).(*utils.DirConfig)
	baseUiUrl := ctx.Value(utils.BaseUiUrl).(string)

	runStatusCode, err := runStatus(httpClient, state)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: err.Error(),
		}
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

	// Following statuses are the statuses you can receive from [creating a brand-new run] to [Run Processing]
	isRunComplete := runStatusCode != models.Ready && runStatusCode != models.Creating && runStatusCode != models.Waiting && runStatusCode != models.Processing
	if !(isRunComplete) {
		outputMsg := fmt.Sprintf("Run number %d - Completed %d out of %d | %s %d, %s %d, %s %d",
			state.RunNumber, len(failedSteps)+len(successSteps), len(steps.Steps),
			goterm.Color("Processing:", goterm.YELLOW), len(processingSteps),
			goterm.Color("Succeeded:", goterm.GREEN), len(successSteps),
			goterm.Color("Failed:", goterm.RED), len(failedSteps))
		if len(processingSteps) != 0 {
			outputMsg += fmt.Sprintf(" | Processing steps: 🥁 %s 🥁", strings.Join(processingSteps, ","))
		}
		if len(failedSteps) != 0 {
			outputMsg += fmt.Sprintf(" | Failed steps: 💩 %s 💩", goterm.Color(strings.Join(failedSteps, ","), goterm.RED))
		}

		return Status{
			PipelinesStatus: "processing",
			Message:         outputMsg,
			Type:            InProgress,
		}
	}

	testsFailureOutput, err := createTestReport(httpClient, state, stepsIdToNames)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps test reports for run id '%d': %v", state.RunId, err),
		}
	}

	outputMsg := fmt.Sprintf("Run %d has completed %d steps with status %s | %s %d, %s %d, %s %d",
		state.RunNumber, len(steps.Steps), runStatusCode.StatusCodeName(),
		goterm.Color("Processing:", goterm.YELLOW), len(processingSteps),
		goterm.Color("Succeeded:", goterm.GREEN), len(successSteps),
		goterm.Color("Failed:", goterm.RED), len(failedSteps))
	if len(failedSteps) != 0 {
		outputMsg += fmt.Sprintf("\n Failed steps: 💩 %s 💩",
			goterm.Color(strings.Join(failedSteps, ","), goterm.RED))
	}
	if testsFailureOutput != "" {
		outputMsg += fmt.Sprintf("\nTests results: %s", testsFailureOutput)
	}

	return Status{
		Message: outputMsg,
		Link: fmt.Sprintf("%s\n",
			utils.GetPipelinesRunURL(baseUiUrl, dirConfig.PipelineName, dirConfig.DefaultStep, state.RunNumber, state.GitBranch)),
		Type: Done,
	}
}

func runStatus(httpClient http.PipelineHttpClient, state *PipedCommandState) (models.StatusCode, error) {
	runRes, err := requests.GetRuns(httpClient, models.GetRunsOptions{
		RunIds: strconv.Itoa(state.RunId),
	})
	if err != nil {
		return models.Failure, fmt.Errorf("failed fetching pipeline runs data: %v", err)
	}
	if len(runRes.Runs) == 0 {
		return models.Failure, fmt.Errorf("failed fetching pipeline runs data: %v", err)
	}

	return runRes.Runs[0].StatusCode, nil
}

func createTestReport(httpClient http.PipelineHttpClient, state *PipedCommandState, stepsIdToNames map[int]string) (string, error) {
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
	return nil
}
