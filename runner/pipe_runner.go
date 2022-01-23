package runner

import (
	"context"
	"github.com/hanochg/piperika/runner/command"
	"github.com/hanochg/piperika/terminal"
	"github.com/pkg/errors"
	"time"
)

var (
	shortBackoffConfig  = backoffConfig{interval: 5 * time.Second, firstTimeout: 15 * time.Second, afterTriggerTimeout: 5 * time.Minute}
	mediumBackoffConfig = backoffConfig{interval: 5 * time.Second, firstTimeout: 30 * time.Second, afterTriggerTimeout: 15 * time.Minute}
	longBackoffConfig   = backoffConfig{interval: 10 * time.Second, firstTimeout: 2 * time.Hour, afterTriggerTimeout: 1 * time.Hour}

	cmds = []pipedCommand{
		newRetryingPipedCommand("Git state", "", command.New001ValidateGitState(), shortBackoffConfig),
		newRetryingPipedCommand("Wait or trigger pipelines sources sync", "sync pipelines sources", command.New002PipelinesSourcesBranchSync(), mediumBackoffConfig),
		newRetryingPipedCommand("Finding active run or trigger", "trigger a run", command.New003PipelinesFindRun(), shortBackoffConfig),
		newRetryingPipedCommand("Wait for run to finish", "", command.New004PipelinesWaitRun(), shortBackoffConfig),
		newRetryingPipedCommand("Run results", "", command.New005PipelinesPrintRun(), longBackoffConfig),
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
			var unRecErr *unrecoverableError
			if errors.As(err, &unRecErr) {
				if termErr := terminal.UpdateUnrecoverable(cmd.OperationName(), unRecErr.Message, ""); termErr != nil {
					return termErr
				} else {
					return nil
				}
			}
			return err
		}
	}

	return nil
}
