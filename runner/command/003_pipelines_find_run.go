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
	"time"
)

func New003PipelinesFindRun() *_003 {
	return &_003{}
}

type _003 struct{}

func (c *_003) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	dirConfig := ctx.Value(utils.DirConfigCtxKey).(*utils.DirConfig)

	pipeResp, err := requests.GetPipelines(httpClient, models.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   state.GitBranch,
		Light:      true,
		PipesNames: dirConfig.PipelineName,
	})
	if err != nil {
		return Status{
			Type:    InProgress,
			Message: fmt.Sprintf("Failed fetching pipelines data: %v", err),
		}
	}
	if len(pipeResp.Pipelines) == 0 {
		return Status{
			PipelinesStatus: "missing pipeline",
			Message:         "Waiting for pipeline creation",
			Type:            InProgress,
		}
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
		return Status{
			Type:    InProgress,
			Message: fmt.Sprintf("Failed fetching pipeline runs data: %v", err),
		}
	}
	if len(runResp.Runs) == 0 {
		// Itai comment: not sure about that, runs handling should be in next step.
		return Status{
			Type:    InProgress,
			Message: "No runs exist for this pipeline branch, triggering new run",
		}
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
		return Status{
			Type:    InProgress,
			Message: fmt.Sprintf("Failed fetching run resources data: %v", err),
		}
	}

	if len(runResourceResp.Resources) == 0 {
		return Status{
			Type:    Failed,
			Message: "No resources exist for the resolved pipeline run, triggering new run",
		}
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
			return Status{
				Type:    Failed,
				Message: "Corrupted data for the resolved pipeline run, triggering new run",
			}
		}
		runNumber, err := strconv.Atoi(runNumbers[i])
		if err != nil {
			return Status{
				Type:    Failed,
				Message: "Corrupted data for the resolved pipeline run, triggering new run",
			}
		}
		if utils.Contains(activeRunIds, runId) {
			state.RunId = runId
			state.RunNumber = runNumber
			break
		}
	}

	if len(activeRunIds) != 0 && state.RunId != -1 {
		return Status{
			Type:    Done,
			Message: "Found an active run id",
		}
	}

	return Status{
		Type:    Failed,
		Message: "did not find any active runs, triggering new run",
	}
}

func (c *_003) TriggerStateChange(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	pipeSteps, err := requests.GetPipelinesSteps(httpClient, models.GetPipelinesStepsOptions{
		PipelineIds:       strconv.Itoa(state.PipelineId),
		PipelineSourceIds: strconv.Itoa(state.PipelinesSourceId),
		Names:             utils.DefaultPipelinesStepNameToTrigger,
	})
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline steps: %v", err),
		}
	}
	if len(pipeSteps.Steps) == 0 {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("No pipeline step called '%s'", utils.DefaultPipelinesStepNameToTrigger),
		}
	}

	err = requests.TriggerPipelinesStep(httpClient, pipeSteps.Steps[0].Id)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed triggering pipeline step '%s': %v", utils.DefaultPipelinesStepNameToTrigger, err),
		}
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
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("Failed fetching pipeline runs: %v", err),
		}
	}

	if len(runResp.Runs) == 0 {
		return Status{
			Type:    Unrecoverable,
			Message: "No runs exist for the pipeline",
		}
	}

	state.RunId = runResp.Runs[0].RunId
	state.RunNumber = runResp.Runs[0].RunNumber
	return Status{
		Type:    Done,
		Message: "Successfully triggered new pipeline run",
	}
}
