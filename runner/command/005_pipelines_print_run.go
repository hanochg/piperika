package command

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
)

func New005PipelinesPrintRun() *_005 {
	return &_005{}
}

type _005 struct{}

func (c *_005) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	// Get steps statuses
	steps, err := requests.GetSteps(httpClient, models.GetStepsOptions{
		// TODO: shouldn't this also have run number? otherwise it's the steps of the last run and not the intended run
		RunIds: strconv.Itoa(state.RunId),
		Limit:  0,
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps data for run id '%d': %v", state.RunId, err),
		}
	}

	failedSteps := make([]string, 0)
	successSteps := make([]string, 0)
	processingSteps := make([]string, 0)
	for _, step := range steps.Steps {
		if step.StatusCode == models.Failure {
			failedSteps = append(failedSteps, step.Name)
		}
		if step.StatusCode == models.Success {
			successSteps = append(successSteps, step.Name)
		}
		if step.StatusCode == models.Waiting ||
			step.StatusCode == models.Processing {
			processingSteps = append(processingSteps, step.Name)
		}
	}

	if len(processingSteps) != 0 {
		return Status{
			PipelinesStatus: "processing",
			Message: fmt.Sprintf("run %d has %d steps. currently %d are processing, %d failed, and %d succeeded",
				state.RunNumber, len(steps.Steps), len(processingSteps), len(failedSteps), len(successSteps)),
			Type: InProgress,
		}
	}

	_, err = requests.GetStepsTestReports(httpClient, models.StepsTestReportsOptions{StepIds: state.RunStepIdsCsv})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps test reports for run id '%d': %v", state.RunId, err),
		}
	}
	// TODO - print the tests results

	return Status{
		Message: fmt.Sprintf("run %d has %d steps. %d failed, and %d succeeded",
			state.RunNumber, len(steps.Steps), len(failedSteps), len(successSteps)),
		Type: Done,
	}
}

func (c *_005) TriggerOnFail(ctx context.Context, state *PipedCommandState) error {
	return fmt.Errorf("timed out")
}
