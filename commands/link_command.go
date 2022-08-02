package commands

import (
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/terminal"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
)

func GetLinkCommand() components.Command {
	return components.Command{
		Name:        "link",
		Description: "Get Pipelines Link",
		Aliases:     []string{"l"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		Action:      getPipelinesLink,
	}
}

func getPipelinesLink(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}

	config, err := utils.GetConfigurations()
	if err != nil {
		return err
	}

	uiUrl, err := utils.GetUIBaseUrl(c)
	if err != nil {
		return err
	}
	branchName, err := utils.GetCurrentBranchName(c)
	if err != nil {
		return err
	}
	projName, err := utils.GetProjectNameForSource(client, config.PipelinesSourceId)
	if err != nil {
		return err
	}
	link := fmt.Sprintf("%s ",
		utils.GetPipelinesBranchURL(uiUrl, config.PipelineName, "", branchName, projName))
	return terminal.DoneMessage("Link", "", link)
}
