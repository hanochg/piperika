package runner

import (
	"context"
	"github.com/hanochg/piperika/runner/commands"
	"time"
)

var (
	defaultBackoffConfig = backoffConfig{interval: time.Second, maxRetries: 30}

	cmds = []PipedCommand{
		NewRetryingPipedCommand("fetch branch", command.New001GetPipeSourceBranch(), defaultBackoffConfig),
		NewRetryingPipedCommand("sync", command.New002WaitPipSourceCompletion(), defaultBackoffConfig),
		NewRetryingPipedCommand("create or grab current run", command.New003GetRun(), defaultBackoffConfig),
		NewRetryingPipedCommand("follow run", command.New004WaitForRun(), defaultBackoffConfig),
	}
)

func RunPipe(ctx context.Context) error {
	pipedState := &command.PipedCommandState{}
	for _, cmd := range cmds {
		err := cmd.Run(ctx, pipedState)
		if err != nil {
			return err
		}
	}
	return nil
}
