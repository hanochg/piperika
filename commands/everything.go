package commands

import (
	"context"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/runner"
	"github.com/hanochg/piperika/runner/command"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"time"
)

func everythingCommand(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()
	return runner.RunPipe(ctx, allPipedCommands(client)...)
}

func allPipedCommands(client http.PipelineHttpClient) []command.PipedCommand {
	return []command.PipedCommand{
		command.NewGitDetailsCommand(),
		command.NewPipelinesSyncStatusCommand(client),
		command.NewTriggerPipelinesSyncCommand(client),
	}
}
