package runner

import (
	"context"
	"github.com/hanochg/piperika/runner/command"
	"github.com/hanochg/piperika/terminal"
	"time"
)

var (
	shortBackoffConfig  = backoffConfig{interval: 5 * time.Second, maxRetries: 12}   // 1 minutes
	mediumBackoffConfig = backoffConfig{interval: 5 * time.Second, maxRetries: 60}   // 5 minutes
	longBackoffConfig   = backoffConfig{interval: 10 * time.Second, maxRetries: 360} // 1 hour

	cmds = []pipedCommand{
		newRetryingPipedCommand("Validate git state", "timeout", command.New001ValidateGitState(), mediumBackoffConfig),
		newRetryingPipedCommand("Wait for syncing pipelines sources", "Sync pipelines sources", command.New002PipelinesSourcesBranchSync(), mediumBackoffConfig),
		newRetryingPipedCommand("Find active run", "Trigger a run", command.New003PipelinesFindRun(), shortBackoffConfig),
		newRetryingPipedCommand("Wait for run to finish", "timeout", command.New004PipelinesWaitRun(), longBackoffConfig),
		newRetryingPipedCommand("Print run results", "timeout", command.New005PipelinesPrintRun(), longBackoffConfig),
	}
)

func RunPipe(ctx context.Context) error {
	pipedState := &command.PipedCommandState{}
	for _, cmd := range cmds {
		err := terminal.StartingRun(cmd.OperationName())
		if err != nil {
			return err
		}

		err = cmd.Run(ctx, pipedState)
		if err != nil {
			return err
		}

	}
	return nil
}
