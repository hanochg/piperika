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
	interval   time.Duration
	maxRetries int
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
	tryOnceErr := c.retryResolveState(ctx, state, &backoff.StopBackOff{})
	if tryOnceErr == nil {
		return nil
	}
	if errors.As(tryOnceErr, &unrecoverableError{}) {
		return tryOnceErr
	}

	err := terminal.UpdateFail(c.operationName, c.failState, "", "")
	if err != nil {
		return err
	}

	err = c.TriggerOnFail(ctx, state)
	if err != nil {
		return err
	}

	return c.retryResolveState(ctx, state, c.newBackoffContext(ctx))
}

func (c *retryingPipedCommand) retryResolveState(ctx context.Context, state *command.PipedCommandState, backoffConfig backoff.BackOff) error {
	return backoff.Retry(
		func() error {
			status := c.ResolveState(ctx, state)

			switch status.Type {
			case command.InProgress:
				err := terminal.UpdateStatus(c.operationName, status.PipelinesStatus, status.Message, status.Link)
				if err != nil {
					return backoff.Permanent(errors.Wrap(&unrecoverableError{}, err.Error()))
				}
				return fmt.Errorf("retrying %s: %s", c.operationName, status.Message)
			case command.Failed:
				err := terminal.UpdateFail(c.operationName, status.PipelinesStatus, status.Message, status.Link)
				if err != nil {
					return backoff.Permanent(errors.Wrap(&unrecoverableError{}, err.Error()))
				}

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

func (c *retryingPipedCommand) newBackoffContext(ctx context.Context) backoff.BackOffContext {
	initialBackoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(c.maxRetries))
	return backoff.WithContext(initialBackoff, ctx)
}
