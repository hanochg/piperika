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
	OperationName() string
}

type backoffConfig struct {
	interval            time.Duration
	firstTimeout        time.Duration
	afterTriggerTimeout time.Duration
}

type retryingPipedCommand struct {
	command.Command
	operationName string
	failState     string
	backoffConfig
}

func newRetryingPipedCommand(operationName, failState string, cmd command.Command, backoffConfig backoffConfig) *retryingPipedCommand {
	return &retryingPipedCommand{
		Command:       cmd,
		backoffConfig: backoffConfig,
		operationName: operationName,
		failState:     failState,
	}
}

func (c *retryingPipedCommand) OperationName() string {
	return c.operationName
}

func (c *retryingPipedCommand) Run(ctx context.Context, state *command.PipedCommandState) error {
	err := c.retryResolveState(ctx, state, c.newBackoffContext(ctx, true))
	if err == nil {
		return nil
	}
	if errors.As(err, &unrecoverableError{}) {
		return err
	}

	terminal.UpdateFail(c.operationName, c.failState, "", "")

	err = c.TriggerOnFail(ctx, state)
	if err != nil {
		_ = terminal.UpdateUnrecoverable(c.operationName, err.Error(), "")

		return err
	}

	return c.retryResolveState(ctx, state, c.newBackoffContext(ctx, false))
}

func (c *retryingPipedCommand) retryResolveState(ctx context.Context, state *command.PipedCommandState, backoffConfig backoff.BackOff) error {
	return backoff.Retry(
		func() error {
			status := c.ResolveState(ctx, state)

			switch status.Type {
			case command.InProgress:
				terminal.UpdateStatus(c.operationName, status.PipelinesStatus, status.Message, status.Link)
				return fmt.Errorf("retrying %s: %s", c.operationName, status.Message)
			case command.Failed:
				terminal.UpdateFail(c.operationName, status.PipelinesStatus, status.Message, status.Link)

				return backoff.Permanent(fmt.Errorf(status.Message))
			case command.Unrecoverable:
				err := terminal.UpdateUnrecoverable(c.operationName, status.Message, status.Link)
				if err != nil {
					return backoff.Permanent(errors.Wrap(&unrecoverableError{}, err.Error()))
				}

				return backoff.Permanent(errors.Wrap(&unrecoverableError{}, status.Message))
			case command.Done:
				err := terminal.DoneMessage(c.operationName, status.Message, status.Link)
				if err != nil {
					return backoff.Permanent(errors.Wrap(&unrecoverableError{}, err.Error()))
				}
				return nil
			default:
				panic("Unexpected command type")
			}
		},
		backoffConfig,
	)
}

type unrecoverableError struct {
	error
}

func (c *retryingPipedCommand) newBackoffContext(ctx context.Context, isFirst bool) backoff.BackOffContext {
	timeout := c.firstTimeout
	if !isFirst {
		timeout = c.afterTriggerTimeout
	}
	maxRetries := uint64(timeout.Nanoseconds() / c.interval.Nanoseconds())

	initialBackoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), maxRetries)
	return backoff.WithContext(initialBackoff, ctx)
}
