package command

import (
	"context"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/http"
)

func NewTriggerPipelinesSyncCommand(client http.PipelineHttpClient) PipedCommand {
	const opName = "trigger pipelines sync"
	initialBackoff := backoff.NewExponentialBackOff()

	return NewRetryingPipedCommand(opName, initialBackoff, func(ctx context.Context, state *PipedCommandState) error {
		if !state.ShouldTriggerPipelinesSync {
			logInfo(opName, "pipelines sources already up to date. skipping sync")
			return nil
		}

		// TODO

		logInfo(opName, "finished")
		return nil
	})
}
