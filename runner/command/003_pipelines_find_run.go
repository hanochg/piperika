package command

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/utils"
	"strconv"
	"strings"
	"time"
)

func New003PipelinesFindRun() *_003 {
	return &_003{}
}

type _003 struct{}

func (c *_003) ResolveState(ctx context.Context, state *PipedCommandState) (Status, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	dirConfig := ctx.Value(utils.DirConfigCtxKey).(*utils.DirConfig)

	pipeResp, err := requests.GetPipelines(httpClient, models.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   state.GitBranch,
		Light:      true,
		PipesNames: dirConfig.PipelineName,
	})
	if err != nil {
		return Status{}, err
	}
	if len(pipeResp.Pipelines) == 0 {
		return Status{
			PipelinesStatus: "missing pipeline",
			Message:         "Waiting for pipeline creation",
			Type:            InProgress,
		}, nil
	}
	state.PipelineId = pipeResp.Pipelines[0].PipelineId

	runResp, err := requests.GetRuns(httpClient, models.GetRunsOptions{
		PipelineIds: strconv.Itoa(state.PipelineId),
		Limit:       10,
		Light:       true,
		StatusCodes: strconv.Itoa(int(models.Processing)),
		SortBy:      "runNumber",
		SortOrder:   -1,
	})
	if err != nil {
		return Status{}, err
	}
	if len(runResp.Runs) == 0 {
		// Itai comment: not sure about that, runs handling should be in next step.
		return Status{}, backoff.Permanent(fmt.Errorf("no runs exist, triggering new run"))
	}

	// TODO: change to struct, and convert to ints here instead of later
	runIds := make([]string, 0)
	runNumbers := make([]string, 0)
	for _, run := range runResp.Runs {
		runIds = append(runIds, strconv.Itoa(run.RunId))
		runNumbers = append(runNumbers, strconv.Itoa(run.RunNumber))
	}
	runResourceResp, err := requests.GetRunResourceVersions(httpClient, models.GetRunResourcesOptions{
		PipelineSourceIds: strconv.Itoa(state.PipelinesSourceId),
		RunIds:            strings.Trim(strings.Join(runIds, ","), "[]"),
		SortBy:            "resourceTypeCode",
		SortOrder:         1,
	})
	if err != nil {
		return Status{}, err
	}

	if len(runResourceResp.Resources) == 0 {
		return Status{}, backoff.Permanent(fmt.Errorf("run has no resources, triggering new run"))
	}

	activeRunIds := make([]int, 0)
	for _, runResource := range runResourceResp.Resources {
		if runResource.ResourceTypeCode != models.GitRepo {
			continue
		}
		if runResource.ResourceVersionContentPropertyBag.CommitSha == state.HeadCommitSha {
			activeRunIds = append(activeRunIds, runResource.RunId)
			break
		}
	}

	// Get the most recent run from the list
	for i, runIdStr := range runIds {
		runId, err := strconv.Atoi(runIdStr)
		if err != nil {
			return Status{}, err
		}
		runNumber, err := strconv.Atoi(runNumbers[i])
		if err != nil {
			return Status{}, err
		}
		if utils.Contains(activeRunIds, runId) {
			state.RunId = runId
			state.RunNumber = runNumber
			break
		}
	}

	if len(activeRunIds) != 0 && state.RunId != -1 {
		return Status{
			Message: "Found an active run id",
			Type:    Done,
		}, nil
	}

	return Status{}, backoff.Permanent(fmt.Errorf("did not find any active runs, triggering new run"))
}

func (c *_003) TriggerStateChange(ctx context.Context, state *PipedCommandState) error {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	pipeSteps, err := requests.GetPipelinesSteps(httpClient, models.GetPipelinesStepsOptions{
		PipelineIds:       strconv.Itoa(state.PipelineId),
		PipelineSourceIds: strconv.Itoa(state.PipelinesSourceId),
		Names:             utils.DefaultPipelinesStepNameToTrigger,
	})

	if err != nil {
		return err
	}
	if len(pipeSteps.Steps) == 0 {
		return fmt.Errorf("tried to trigger a run for step '%s' but coulnd't fetch its Id", utils.DefaultPipelinesStepNameToTrigger)
	}

	err = requests.TriggerPipelinesStep(httpClient, pipeSteps.Steps[0].Id)
	if err != nil {
		return err
	}

	// Giving Pipelines time to digest the request and create a new run
	time.Sleep(3 * time.Second)

	runResp, err := requests.GetRuns(httpClient, models.GetRunsOptions{
		PipelineIds: strconv.Itoa(state.PipelineId),
		Limit:       1,
		Light:       true,
		SortBy:      "createdAt",
		SortOrder:   -1,
	})
	if err != nil {
		return err
	}

	if len(runResp.Runs) == 0 {
		return fmt.Errorf("failed to get the triggered run")
	}

	state.RunId = runResp.Runs[0].RunId
	state.RunNumber = runResp.Runs[0].RunNumber
	return nil
}
