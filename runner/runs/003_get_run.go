package runs

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/http/models"
	"github.com/hanochg/piperika/http/requests"
	"github.com/hanochg/piperika/runner/datastruct"
	"github.com/hanochg/piperika/utils"
	"strconv"
)

/*
	Run Description
	---------------
	- Get all the relevant runs
	- Check if there are relevant (on the same commit sha) active runs
	- Trigger a new run with "trigger_all"
*/

func (_ _03) Init(ctx context.Context, state *datastruct.PipedCommandState) (string, error) {
	// TODO set it according the local dir
	state.PipelineName = "access_build"
	return "", nil
}

func (_ _03) Tick(ctx context.Context, state *datastruct.PipedCommandState) (*datastruct.RunStatus, error) {
	httpClient := ctx.Value("client").(http.PipelineHttpClient)
	pipeResp, err := requests.GetPipelines(httpClient, models.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   state.GitBranch,
		Light:      true,
		PipesNames: state.PipelineName,
	})
	if err != nil {
		return nil, err
	}
	if len(pipeResp.Pipelines) == 0 {
		return nil, fmt.Errorf("failed to get the pipeline")
	}
	state.PipelineId = pipeResp.Pipelines[0].PipelineId

	runResp, err := requests.GetRuns(httpClient, models.GetRunsOptions{
		PipelineIds: strconv.Itoa(state.PipelineId),
		Limit:       10,
		Light:       true,
		StatusCodes: models.Processing,
		SortBy:      "runNumber",
		SortOrder:   -1,
	})
	if err != nil {
		return nil, err
	}
	if len(runResp.Runs) == 0 {
		state.ShouldTriggerRun = true
		return &datastruct.RunStatus{
			Done: true,
		}, nil
	}

	activeRunId := -1
	for _, run := range runResp.Runs {
		runEndAtTime, err := utils.PipelinesTimeParser(run.EndedAt)
		if err != nil {
			// TODO log
			continue
		}
		if runEndAtTime.After(state.PipelinesSyncDate) {
			activeRunId = run.RunId
			break
		}
	}

	if activeRunId != -1 {
		state.RunId = activeRunId
		return &datastruct.RunStatus{
			Message: "Found an active run id",
			Done:    true,
		}, nil

	}

	state.ShouldTriggerRun = true
	return &datastruct.RunStatus{
		Message: "No active run found, triggering run",
		Done:    true,
	}, nil
}

func (_ _03) OnComplete(ctx context.Context, state *datastruct.PipedCommandState, status *datastruct.RunStatus) (string, error) {
	return "", nil
}

type _03 struct {
}

func New003GetRun() datastruct.Runner {
	return _03{}
}
