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

func (c *_004) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	// TODO - as the run steps don't change, there's no need ot get them in each cycle
	// we need to separate GetSteps from the GetRuns cycle
	stepsResp, err := requests.GetSteps(httpClient, models.GetStepsOptions{
		// TODO: shouldn't this also have run number? otherwise it's the steps of the last run and not the intended run
		// Hanoch: run id is points to a specific run number in a specific pipeline (e.g access_build under RT repo),
		// each run gets a run id, I added the run number just for printing it to the console (as it's user friendly)
		//
		RunIds: strconv.Itoa(state.RunId),
		Limit:  100,
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps data for run id '%d': %v", state.RunId, err),
		}
	}

	if len(stepsResp.Steps) == 0 {
		// Itai comment - Do we want to wait here or fail?
		return Status{
			Type:    Failed,
			Message: fmt.Sprintf("No steps for pipeline run id '%d'", state.RunId),
		}
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
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline runs data: %v", err),
		}
	}
	if len(runStatus.Runs) == 0 {
		return Status{
			Type:            InProgress,
			PipelinesStatus: "waiting for run",
			Message:         fmt.Sprintf("Could not resolve any runs for the pipeline"),
		}
	}

	statusCode := runStatus.Runs[0].StatusCode
	if statusCode == models.Creating || statusCode == models.Waiting || statusCode == models.Processing {
		return Status{
			PipelinesStatus: models.StatusCodeNamesMap[statusCode],
			Message: fmt.Sprintf("run %d started at %s is in progress",
				runStatus.Runs[0].RunNumber, runStatus.Runs[0].StartedAt),
			Type: InProgress,
		}
	}

	return Status{
		Message: fmt.Sprintf("run %d started at %s and finished at %s with status %d",
			runStatus.Runs[0].RunNumber, runStatus.Runs[0].StartedAt, runStatus.Runs[0].EndedAt, statusCode),
		Type: Done,
	}
}

func (c *_004) TriggerOnFail(ctx context.Context, state *PipedCommandState) error {
	return fmt.Errorf("timed out")
}
