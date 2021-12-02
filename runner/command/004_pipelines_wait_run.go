package command

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
	"strings"
)

func New004PipelinesWaitRun() *_004 {
	return &_004{}
}

type _004 struct{}

func (c *_004) ResolveState(ctx context.Context, state *PipedCommandState) (Status, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	stepsResp, err := requests.GetSteps(httpClient, models.GetStepsOptions{
		// TODO: shouldn't this also have run number? otherwise it's the steps of the last run and not the intended run
		RunIds: strconv.Itoa(state.RunId),
		Limit:  100,
	})
	if err != nil {
		return Status{}, err
	}

	if len(stepsResp.Steps) == 0 {
		return Status{}, fmt.Errorf("cannot get the steps of the current run, run id %d", state.RunId)
	}

	stepIds := make([]string, 0)
	for _, step := range stepsResp.Steps {
		stepIds = append(stepIds, strconv.Itoa(step.Id))
	}
	state.RunStepIdsCsv = strings.Trim(strings.Join(stepIds, ","), "[]")

	runStatus, err := requests.GetRuns(httpClient, models.GetRunsOptions{
		RunIds: strconv.Itoa(state.RunId),
	})
	if err != nil {
		return Status{}, err
	}
	if len(runStatus.Runs) == 0 {
		return Status{}, fmt.Errorf("cannot fetch run, retrying")
	}

	if runStatus.Runs[0].StatusCode == models.Creating ||
		runStatus.Runs[0].StatusCode == models.Waiting ||
		runStatus.Runs[0].StatusCode == models.Processing {
		return Status{}, fmt.Errorf("run %d started at %s is in progress",
			runStatus.Runs[0].RunNumber, runStatus.Runs[0].StartedAt)
	}

	return Status{
		Message: fmt.Sprintf("run %d started at %s and finished at %s with status %d",
			runStatus.Runs[0].RunNumber, runStatus.Runs[0].StartedAt, runStatus.Runs[0].EndedAt, runStatus.Runs[0].StatusCode),
	}, nil
}

func (c *_004) TriggerStateChange(ctx context.Context, state *PipedCommandState) error {
	// do nothing
	return nil
}
