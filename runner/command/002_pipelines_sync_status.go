package command

import (
	"context"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/http"
)

func NewPipelinesSyncStatusCommand(client http.PipelineHttpClient) PipedCommand {
	const opName = "pipelines sync status"
	initialBackoff := backoff.NewExponentialBackOff()

	return NewRetryingPipedCommand(opName, initialBackoff, func(ctx context.Context, state *PipedCommandState) error {
		// TODO: check sync status and compare commit hash to local
		state.ShouldTriggerPipelinesSync = true
		logInfo(opName, "finished")
		return nil
	})
}
