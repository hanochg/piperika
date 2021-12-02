package commands

import (
	"context"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/runner"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"time"
)

func GetCommand() components.Command {
	return components.Command{
		Name:        "pipelines",
		Description: "pipelines operations",
		Aliases:     []string{"pp"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		Action:      action,
	}
}

func getArguments() []components.Argument {
	return []components.Argument{
		{Name: "command", Description: "perform a command like get / list"},
	}
}

func getFlags() []components.Flag {
	return []components.Flag{plugins.GetServerIdFlag()}
}

func action(c *components.Context) error {
	return theAllMightyCommand(c)
}

func theAllMightyCommand(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	//dirConfig, err := utils.GetDirConfig()
	dirConfig := &utils.DirConfig{ //TODO for tests
		PipelineName:      "access_build",
		PipelinesSourceId: 6,
	}
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()
	ctx = context.WithValue(ctx, utils.HttpClientCtxKey, client)
	ctx = context.WithValue(ctx, utils.DirConfigCtxKey, dirConfig)
	return runner.RunPipe(ctx)
}
