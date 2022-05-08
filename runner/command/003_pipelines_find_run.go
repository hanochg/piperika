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

func New003PipelinesFindRun() *_003 {
	return &_003{}
}

type _003 struct {
	runTriggered bool
}

func (c *_003) ResolveState(ctx context.Context, state *PipedCommandState) Status {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	dirConfig := ctx.Value(utils.ConfigCtxKey).(*utils.Configurations)
	branchName := ctx.Value(utils.BranchName).(string)
	forceFlag := ctx.Value(utils.ForceFlag).(bool)

	pipelineId, err := getPipelineIdByBranch(httpClient, dirConfig.PipelineName, branchName)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: err.Error(),
		}
	}
	if pipelineId == -1 {
		return Status{
			Type:            InProgress,
			PipelinesStatus: "missing pipeline",
			Message:         fmt.Sprintf("waiting for pipeline '%s' creation", dirConfig.PipelineName),
		}
	}
	state.PipelineId = pipelineId

	runs, err := getRuns(httpClient, pipelineId, forceFlag)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: fmt.Sprintf("failed fetching pipeline runs data: %v", err),
		}
	}
	if len(runs) == 0 {
		status := "there are no runs"
		if forceFlag {
			status = "there are no processing runs"
		}
		return Status{
			PipelinesStatus: status,
			Message:         "waiting for run creation",
			Type:            InProgress,
		}
	}

	activeRunId, err := getActiveRunId(httpClient, state.PipelinesSourceId, runs, state.HeadCommitSha)
	if err != nil {
		return Status{
			Type:    Unrecoverable,
			Message: err.Error(),
		}
	}

	if activeRunId == -1 {
		return Status{
			Type:            InProgress,
			PipelinesStatus: "Not exists",
			Message:         "did not find any active runs",
		}
	}

	state.RunId = activeRunId
	state.RunNumber = getRunNumberById(runs, activeRunId)

	if c.runTriggered {
		return Status{
			Message: fmt.Sprintf("Run #%d was triggered", state.RunNumber),
			Type:    Done,
		}
	}
	return Status{
		Message: fmt.Sprintf("Run #%d was found", state.RunNumber),
		Type:    Done,
	}
}

func getRuns(httpClient http.PipelineHttpClient, pipelineId int, onlyRunning bool) ([]requests.Run, error) {
	statusCode := ""
	if onlyRunning {
		statusCode = strings.Join([]string{http.Ready.String(), http.Creating.String(), http.Waiting.String(), http.Processing.String()}, ",")
	}
	runs, err := requests.GetRuns(httpClient, requests.GetRunsOptions{
		PipelineIds: strconv.Itoa(pipelineId),
		Limit:       100,
		Light:       true,
		StatusCodes: statusCode,
		SortBy:      "runNumber",
		SortOrder:   -1,
	})
	if err != nil {
		return nil, err
	}
	return runs.Runs, nil
}

func getRunNumberById(runs []requests.Run, runId int) int {
	for _, run := range runs {
		if run.RunId == runId {
			return run.RunNumber
		}
	}
	return -1
}

func getActiveRunId(httpClient http.PipelineHttpClient, pipelineSourceId int, runs []requests.Run, headCommitHash string) (int, error) {
	runIds := make([]string, 0)
	for _, run := range runs {
		runIds = append(runIds, strconv.Itoa(run.RunId))
	}
	runIdsList := strings.Trim(strings.Join(runIds, ","), "[]")
	runResourceResp, err := requests.GetRunResourceVersions(httpClient, requests.GetRunResourcesOptions{
		Limit:             10000,
		PipelineSourceIds: strconv.Itoa(pipelineSourceId),
		RunIds:            runIdsList,
		SortBy:            "createdAt",
		SortOrder:         -1,
	})
	if err != nil {
		return 0, fmt.Errorf("failed fetching run resources data: %w", err)
	}

	activeRunIds := getActiveRunIdsByResources(runResourceResp.Resources, headCommitHash)
	return activeRunIds, nil
}

func getActiveRunIdsByResources(resources []requests.RunResource, headCommit string) int {
	for _, runResource := range resources {
		if runResource.ResourceTypeCode == http.GitRepo &&
			runResource.ResourceVersionContentPropertyBag.CommitSha == headCommit {
			return runResource.RunId
		}
	}
	return -1
}

func getPipelineIdByBranch(client http.PipelineHttpClient, pipelineName, branchName string) (int, error) {
	pipeResp, err := requests.GetPipelines(client, requests.GetPipelinesOptions{
		SortBy:     "latestRunId",
		FilterBy:   branchName,
		Light:      true,
		PipesNames: pipelineName,
	})
	if err != nil {
		return 0, fmt.Errorf("failed fetching pipelines data: %w", err)
	}

	for _, pipeline := range pipeResp.Pipelines {
		if pipeline.PipelineSourceBranch == branchName {
			return pipeline.PipelineId, nil
		}
	}

	return -1, nil
}

func (c *_003) TriggerOnFail(ctx context.Context, state *PipedCommandState) error {
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	dirConfig := ctx.Value(utils.ConfigCtxKey).(*utils.Configurations)

	pipeSteps, err := requests.GetPipelinesSteps(httpClient, requests.GetPipelinesStepsOptions{
		PipelineIds:       strconv.Itoa(state.PipelineId),
		PipelineSourceIds: strconv.Itoa(state.PipelinesSourceId),
		Names:             dirConfig.DefaultStep,
	})
	if err != nil {
		return fmt.Errorf("failed fetching pipeline steps: %w", err)
	}
	if len(pipeSteps.Steps) == 0 {
		return fmt.Errorf("no pipeline step called '%s'", dirConfig.DefaultStep)
	}
	c.runTriggered = true
	return requests.TriggerPipelinesStep(httpClient, pipeSteps.Steps[0].Id)
}
