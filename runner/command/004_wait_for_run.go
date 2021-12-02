package command

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/runner/datastruct"
	"github.com/hanochg/piperika/utils"
	"strconv"
	"strings"
)

/*
	Run Description
	---------------
	- Get the run steps and status
	- ResolveState for the run to complete of fail
	- Gives statistics and details about the current run
*/

func (_ _04) RetryableDoBeforeTrigger(ctx context.Context, state *datastruct.PipedCommandState) (string, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	stepsResp, err := requests.GetSteps(httpClient, models.GetStepsOptions{
		RunIds: strconv.Itoa(state.RunId),
		Limit:  100,
	})
	if err != nil {
		return "", err
	}

	if len(stepsResp.Steps) == 0 {
		return "", fmt.Errorf("cannot get the steps of the current run, run id %d", state.RunId)
	}

	stepIds := make([]string, 0)
	for _, step := range stepsResp.Steps {
		stepIds = append(stepIds, strconv.Itoa(step.Id))
	}
	state.RunStepIdsCsv = strings.Trim(strings.Join(stepIds, ","), "[]")
	return "", nil
}

func (_ _04) TriggerOnceIfNecessary(ctx context.Context, state *datastruct.PipedCommandState) (*datastruct.RunStatus, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	runStatus, err := requests.GetRuns(httpClient, models.GetRunsOptions{
		RunIds: strconv.Itoa(state.RunId),
	})
	if err != nil {
		return nil, err
	}
	if len(runStatus.Runs) == 0 {
		return &datastruct.RunStatus{
			Status: "cannot fetch run, retrying...",
			Done:   false,
		}, nil
	}

	if runStatus.Runs[0].StatusCode != models.Creating &&
		runStatus.Runs[0].StatusCode != models.Waiting &&
		runStatus.Runs[0].StatusCode != models.Processing {
		return &datastruct.RunStatus{
			Status: fmt.Sprintf("run %d started at %s and finished at %s with status %d",
				runStatus.Runs[0].RunNumber, runStatus.Runs[0].StartedAt, runStatus.Runs[0].EndedAt, runStatus.Runs[0].StatusCode),
			Done: true,
		}, nil
	}

	// Get steps statuses
	steps, err := requests.GetSteps(httpClient, models.GetStepsOptions{
		RunIds: strconv.Itoa(state.RunId),
		Limit:  0,
	})
	if err != nil {
		return nil, err
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

	return &datastruct.RunStatus{
		Status: fmt.Sprintf("run %d started at %s and is still running...", runStatus.Runs[0].RunNumber, runStatus.Runs[0].StartedAt),
		Message: fmt.Sprintf("run %d has %d steps. currently %d are processing, %d failed, and %d succeed",
			runStatus.Runs[0].RunNumber, len(steps.Steps), len(processingSteps), len(failedSteps), len(successSteps)),
		Done: true,
	}, nil
}

func (_ _04) RetryableDoAfterTrigger(ctx context.Context, state *datastruct.PipedCommandState, status *datastruct.RunStatus) (string, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	_, err := requests.GetStepsTestReports(httpClient, models.StepsTestReportsOptions{StepIds: state.RunStepIdsCsv})
	if err != nil {
		return "", err

	}
	// TODO - print the tests results

	return "", nil
}

type _04 struct {
}

func New004WaitForRun() datastruct.Runner {
	return _04{}
}
