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
	"strings"
)

/*
	Run Description
	---------------
	- Get all the relevant runs
	- Check if there are relevant (on the same commit sha) active runs
	- Trigger a new run with "trigger_all"
*/

func (_ _03) Init(ctx context.Context, state *datastruct.PipedCommandState) (string, error) {
	return "", nil
}

func (_ _03) Tick(ctx context.Context, state *datastruct.PipedCommandState) (*datastruct.RunStatus, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	dirConfig := ctx.Value(utils.DirConfigCtxKey).(*utils.DirConfig)

	pipeResp, err := requests.GetPipelines(httpClient, models.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   state.GitBranch,
		Light:      true,
		PipesNames: dirConfig.PipelineName,
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

	runIds := make([]string, 0)
	for _, run := range runResp.Runs {
		runIds = append(runIds, strconv.Itoa(run.RunId))
	}
	runResourceResp, err := requests.GetResourceVersions(httpClient, models.GetResourcesOptions{
		PipelineSourceIds: strconv.Itoa(state.PipelinesSourceId),
		RunIds:            strings.Trim(strings.Join(runIds, ","), "[]"),
	})
	if err != nil {
		return nil, err
	}

	if len(runResourceResp.Resources) == 0 {
		state.ShouldTriggerRun = true
		return &datastruct.RunStatus{
			Done: true,
		}, nil
	}

	activeRunId := -1
	for _, runResource := range runResourceResp.Resources {
		if runResource.ResourceTypeCode != models.GitRepo {
			continue
		}
		if runResource.ResourceVersionContentPropertyBag.CommitSha == state.HeadCommitSha {
			activeRunId = runResource.Id
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
		Message: "Triggering a new run",
		Done:    true,
	}, nil

}

func (_ _03) OnComplete(ctx context.Context, state *datastruct.PipedCommandState, status *datastruct.RunStatus) (string, error) {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)

	if state.ShouldTriggerRun {
		pipeSteps, err := requests.GetPipelinesSteps(httpClient, models.GetStepsOptions{
			PipelineIds:       strconv.Itoa(state.PipelineId),
			PipelineSourceIds: strconv.Itoa(state.PipelinesSourceId),
			Names:             utils.DefaultPipelinesStepNameToTrigger,
		})

		if err != nil {
			return "", err
		}
		if len(pipeSteps.Steps) == 0 {
			return "", fmt.Errorf("tried to trigger a run for step '%s' but coulnd't fetch its Id", utils.DefaultPipelinesStepNameToTrigger)
		}

		_, err = requests.TriggerPipelinesStep(httpClient, models.GetStepsOptions{}, pipeSteps.Steps[0].Id)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

type _03 struct {
}

func New003GetRun() datastruct.Runner {
	return _03{}
}
