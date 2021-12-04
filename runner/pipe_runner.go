package runner

import (
	"context"
	"github.com/hanochg/piperika/runner/command"
	"time"
)

var (
	shortBackoffConfig = backoffConfig{interval: 5 * time.Second, maxRetries: 60}   // 5 minutes
	longBackoffConfig  = backoffConfig{interval: 10 * time.Second, maxRetries: 360} // 1 hour

	cmds = []pipedCommand{
		newRetryingPipedCommand("validate git state", command.New001ValidateGitState(), shortBackoffConfig),
		newRetryingPipedCommand("sync pipelines sources", command.New002PipelinesSourcesBranchSync(), shortBackoffConfig),
		newRetryingPipedCommand("find or trigger active run", command.New003PipelinesFindRun(), shortBackoffConfig),
		newRetryingPipedCommand("wait for run to finish", command.New004PipelinesWaitRun(), longBackoffConfig),
		newRetryingPipedCommand("print run results", command.New005PipelinesPrintRun(), shortBackoffConfig),
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
