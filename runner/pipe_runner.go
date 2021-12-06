package runner

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/runner/command"
	"github.com/hanochg/piperika/terminal"
	"time"
)

var (
	mediumBackoffConfig = backoffConfig{interval: 5 * time.Second, maxRetries: 60}   // 5 minutes
	longBackoffConfig   = backoffConfig{interval: 10 * time.Second, maxRetries: 360} // 1 hour

	cmds = []pipedCommand{
		newRetryingPipedCommand("Git state", "timeout", command.New001ValidateGitState(), mediumBackoffConfig),
		newRetryingPipedCommand("Wait or trigger pipelines sources sync", "sync pipelines sources", command.New002PipelinesSourcesBranchSync(), mediumBackoffConfig),
		newRetryingPipedCommand("Finding active run or trigger", "trigger a run", command.New003PipelinesFindRun(), mediumBackoffConfig),
		newRetryingPipedCommand("Wait for run to finish", "timeout", command.New004PipelinesWaitRun(), longBackoffConfig),
		newRetryingPipedCommand("Run results", "timeout", command.New005PipelinesPrintRun(), longBackoffConfig),
	}
)

func RunPipe(ctx context.Context) error {
	pipedState := &command.PipedCommandState{}
	for _, cmd := range cmds {
		err := terminal.StartingRun(cmd.OperationName())
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}

		err = cmd.Run(ctx, pipedState)
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}

	}
	return nil
}
