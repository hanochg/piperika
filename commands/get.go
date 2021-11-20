package commands

import (
	"fmt"
	"github.com/hanochg/piperika/steps"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
)

func getCommand(c *components.Context) error {
	switch c.Arguments[1] { // Consider module arguments instead of taking the 2nd argument
	case "step": // Consider use strategy / map instead
		return getRunningStepsForBranch(c)
	}
	return nil
}

func getRunningStepsForBranch(c *components.Context) error {
	client, err := utils.NewPipelineHttp(c)
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
