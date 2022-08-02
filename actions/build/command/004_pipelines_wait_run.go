package command

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
	"strings"
)

func New004PipelinesWaitRun() *_004 {
	return &_004{}
}

type _004 struct{}

func (c *_004) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	stepsResp, err := requests.GetSteps(httpClient, requests.GetStepsOptions{
		RunIds: strconv.Itoa(state.RunId),
		Limit:  50,
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps data for run id '%d': %v", state.RunId, err),
		}
	}

	if len(stepsResp.Steps) == 0 {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("No steps for pipeline run id '%d'", state.RunId),
		}
	}

	stepIds := make([]string, 0)
	for _, step := range stepsResp.Steps {
		stepIds = append(stepIds, strconv.Itoa(step.Id))
	}
	state.RunStepIdsCsv = strings.Trim(strings.Join(stepIds, ","), "[]")

	runResp, err := requests.GetRuns(httpClient, requests.GetRunsOptions{
		RunIds: strconv.Itoa(state.RunId),
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline runs data: %v", err),
		}
	}
	if len(runResp.Runs) == 0 {
		return Status{
			Type:            InProgress,
			PipelinesStatus: "waiting for run",
			Message:         fmt.Sprintf("Could not resolve any runs for the pipeline"),
		}
	}

	return Status{
		Message: fmt.Sprintf("Watching run #%d with %d steps, current run status %s",
			state.RunNumber, len(stepsResp.Steps), runResp.Runs[0].StatusCode.StatusCodeName()),
		Type: Done,
	}
}

func (c *_004) TriggerOnFail(_ context.Context, _ *PipedCommandState) error {
	return nil
}
