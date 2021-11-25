package commands

import (
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/steps"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
)

func getCommand(c *components.Context) error {
	switch c.Arguments[1] { // Consider module arguments instead of taking the 2nd argument
	case "step": // Consider use strategy / map instead
		return getRunningStepsForBranch(c)
	case "pipelines":
		return getPipelinesForBranch(c)
	case "sync":
		return getSyncStatusForBranch(c)
	case "source":
		return getSourcesById(c)
	case "step_connection":
		return getStepConnectionById(c)
	case "runs":
		return getRuns(c)
	case "step_test_reports":
		return getStepTestReports(c)
	case "sync_sources":
		return syncSourcesById(c)

	}
	return nil
}

func getRunningStepsForBranch(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	branch, err := utils.GetCurrentBranchName()
	if err != nil {
		return err
	}
	res, err := steps.GetRunningStepsForBranch(client, branch)
	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}

func getPipelinesForBranch(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	branch, err := utils.GetCurrentBranchName()
	if err != nil {
		return err
	}

	res, err := steps.GetPipelinesForBranch(client, branch)

	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}

func getSyncStatusForBranch(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	branch, err := utils.GetCurrentBranchName()
	if err != nil {
		return err
	}

	res, err := steps.GetSyncStatusForBranch(client, 6, branch)

	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}

func getSourcesById(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}

	res, err := steps.GetSourcesById(client, "6") // TODO input of this

	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}

func syncSourcesById(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	branch, err := utils.GetCurrentBranchName()
	if err != nil {
		return err
	}

	res, err := steps.SyncSources(client, 6, branch) // TODO input of this

	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}

func getStepConnectionById(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}

	res, err := steps.GetStepConnectionsByPipelinesId(client, "118524") // TODO input of this

	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}

func getRuns(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	res, err := steps.GetRuns(client, "269978") // TODO input of this

	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}

func getStepTestReports(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	res, err := steps.GetStepsTestReports(client, "3306611") // TODO input of this

	if err != nil {
		return err
	}
	fmt.Println(res) // TODO Improve this
	return err
}
