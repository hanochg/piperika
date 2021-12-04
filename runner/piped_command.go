package runner

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/runner/command"
	"github.com/hanochg/piperika/terminal"
	"github.com/pkg/errors"
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
	if errors.As(waitErr, &unrecoverableError{}) {
		return waitErr
	}

	status := c.TriggerStateChange(ctx, state)
	if status.Type == command.Unrecoverable {
		return fmt.Errorf(status.Message)
	}
	err := terminal.UpdateStatus(c.operationName, status.PipelinesStatus, status.Message, "TBD", false)
	if err != nil {
		return err
	}

	// Giving Pipelines time to digest the triggered request
	time.Sleep(3 * time.Second)

	return c.retryResolveState(ctx, state)
}

func (c *retryingPipedCommand) retryResolveState(ctx context.Context, state *command.PipedCommandState) error {
	return backoff.Retry(
		func() error {
			status := c.ResolveState(ctx, state)

			isTempLine := status.Type == command.InProgress
			err := terminal.UpdateStatus(c.operationName, status.PipelinesStatus, status.Message, "TBD", isTempLine)
			if err != nil {
				return err
			}

			if status.Type == command.InProgress {
				return fmt.Errorf("retrying %s", c.operationName)
			}
			if status.Type == command.Failed {
				return backoff.Permanent(fmt.Errorf(status.Message))
			}
			if status.Type == command.Unrecoverable {
				return backoff.Permanent(errors.Wrap(&unrecoverableError{}, status.Message))
			}

			// Done
			return nil
		},
		c.newBackoffContext(ctx),
	)
}

type unrecoverableError struct {
	error
}

func (c *retryingPipedCommand) newBackoffContext(ctx context.Context) backoff.BackOffContext {
	initialBackoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(c.maxRetries))
	return backoff.WithContext(initialBackoff, ctx)
}
