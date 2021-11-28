package command

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	"time"
)

// state should only be altered in case of successful operation
type pipedCommandFunc func(ctx context.Context, state *PipedCommandState) error

func NewRetryingPipedCommand(operationName string, initialBackoff backoff.BackOff, operation pipedCommandFunc) *retryingPipedCommand {
	return &retryingPipedCommand{
		operationName:  operationName,
		operation:      operation,
		initialBackoff: initialBackoff,
	}
}

type retryingPipedCommand struct {
	operationName  string
	operation      pipedCommandFunc
	initialBackoff backoff.BackOff
}

func (c *retryingPipedCommand) Run(ctx context.Context, state *PipedCommandState) CommandResolution {
	backoffConfig := backoff.WithContext(c.initialBackoff, ctx)

	err := backoff.RetryNotify(
		func() error {
			return c.operation(ctx, state)
		},
		backoffConfig,
		func(err error, duration time.Duration) {
			logWarning(c.operationName, fmt.Sprintf("%s; retrying in %d", err.Error(), duration))
		},
	)

	return CommandResolution{
		Error: err,
	}
}
