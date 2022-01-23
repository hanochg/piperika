package commands

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/runner"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"time"
)

func GetCommand() components.Command {
	return components.Command{
		Name:        "run",
		Description: "Start a Pipelines cycle with your Git commit and branch",
		Aliases:     []string{"r"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		Action:      action,
	}
}

func getArguments() []components.Argument {
	return []components.Argument{}
}

func getFlags() []components.Flag {
	return []components.Flag{
		plugins.GetServerIdFlag(),
		components.StringFlag{
			Name:        "branch",
			Description: "Specify the branch to build",
		},
		components.BoolFlag{
			Name:        "force",
			Description: "Force trigger if there is no processing runs",
		},
	}
}

func action(c *components.Context) error {
	return theAllMightyCommand(c)
}

func theAllMightyCommand(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	dirConfig, err := utils.GetDirConfig()
	if err != nil {
		return err
	}

	uiUrl, err := utils.GetUIBaseUrl(c)
	if err != nil {
		return err
	}
	branch, err := utils.GetCurrentBranchName(c)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()
	ctx = context.WithValue(ctx, utils.BranchName, branch)
	ctx = context.WithValue(ctx, utils.BaseUiUrl, uiUrl)
	ctx = context.WithValue(ctx, utils.HttpClientCtxKey, client)
	ctx = context.WithValue(ctx, utils.DirConfigCtxKey, dirConfig)
	ctx = context.WithValue(ctx, utils.ForceFlag, c.GetBoolFlagValue("force"))
	err = runner.RunPipe(ctx)
	fmt.Println()
	return err
}
