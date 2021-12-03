package runner

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/runner/command"
	"github.com/hanochg/piperika/terminal"
	"time"
)

type pipedCommand interface {
	Run(ctx context.Context, state *command.PipedCommandState) error
}

type backoffConfig struct {
	interval   time.Duration
	maxRetries int
}

type retryingPipedCommand struct {
	command.Command
	operationName string
	backoffConfig
}

func newRetryingPipedCommand(operationName string, cmd command.Command, backoffConfig backoffConfig) *retryingPipedCommand {
	return &retryingPipedCommand{
		Command:       cmd,
		backoffConfig: backoffConfig,
		operationName: operationName,
	}
}

func (c *retryingPipedCommand) Run(ctx context.Context, state *command.PipedCommandState) error {
	waitErr := c.retryResolveState(ctx, state)

	if waitErr == nil {
		return nil
	}
	if !errors.As(waitErr, &timeOutError{}) {
		return waitErr
	}

	// Time out
	triggerErr := c.TriggerStateChange(ctx, state)
	if triggerErr != nil {
		return triggerErr
	}

	return c.retryResolveState(ctx, state)
}

func (c *retryingPipedCommand) retryResolveState(ctx context.Context, state *command.PipedCommandState) error {
	backoffConfig := c.newBackoffContext(ctx)
	for currInterval := backoffConfig.NextBackOff(); currInterval > 0; currInterval = backoffConfig.NextBackOff() {
		var currentStatus *command.Status
		status, err := c.ResolveState(ctx, state)
		if err != nil {
			return err
		}

		err = terminal.UpdateStatus(c.operationName, status.PipelinesStatus, status.Message, "TBD")
		if err != nil {
			return err
		}

		if currentStatus.Type == command.Done {
			return nil
		}
	}

	return fmt.Errorf("timed-out %w", timeOutError{})
}

type timeOutError struct {
	error
}

func (c *retryingPipedCommand) newBackoffContext(ctx context.Context) backoff.BackOffContext {
	initialBackoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(c.maxRetries))
	return backoff.WithContext(initialBackoff, ctx)
}
