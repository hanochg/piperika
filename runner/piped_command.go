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
	lastStatus, err := c.resolveCurrentState(ctx, state, c.newBackoffContext(ctx, true))
	if err == nil {
		return nil
	}
	var unrecErr *unrecoverableError
	if errors.As(err, &unrecErr) {
		return err
	}

	terminal.UpdateFail(c.operationName, c.failState, "time-out", "")

	err = c.TriggerOnFail(ctx, state)
	if err != nil {
		return terminal.UpdateUnrecoverable(c.operationName, err.Error(), "")
	}

	lastStatus, err = c.resolveCurrentState(ctx, state, c.newBackoffContext(ctx, false))
	if err == nil {
		return nil
	}
	if errors.As(err, &unrecErr) {
		return err
	}

	if termErr := terminal.UpdateUnrecoverable(c.operationName, "timed-out after trigger once", lastStatus.Link); termErr != nil {
		return termErr
	}
	return err
}

func (c *retryingPipedCommand) resolveCurrentState(ctx context.Context, state *command.PipedCommandState, backoffConfig backoff.BackOff) (command.Status, error) {
	lastStatus := command.Status{}
	status := command.Status{}
	for currInterval := time.Nanosecond; currInterval != backoff.Stop; currInterval = backoffConfig.NextBackOff() {
		select {
		case <-time.Tick(currInterval):
			status = c.ResolveState(ctx, state)

			switch status.Type {
			case command.InProgress:
				if lastStatus.PipelinesStatus != status.PipelinesStatus && lastStatus.Message != status.Message {
					backoffConfig.Reset()
					terminal.UpdateStatus(c.operationName, status.PipelinesStatus, status.Message, status.Link)
				}
			case command.Unrecoverable:
				return command.Status{}, &unrecoverableError{
					Status:  status.PipelinesStatus,
					Message: status.Message,
				}
			case command.Done:
				err := terminal.DoneMessage(c.operationName, status.Message, status.Link)
				if err != nil {
					return command.Status{}, err
				}
				return status, nil
			default:
				panic("Unexpected command type")
			}
		case <-ctx.Done():
			return status, nil
		}

		lastStatus = status
	}

	return status, timeOutError{}
}

type timeOutError struct {
}

func (t timeOutError) Error() string {
	return "Time-out"
}

type unrecoverableError struct {
	Status  string
	Message string
}

func (u *unrecoverableError) Error() string {
	return fmt.Sprintf("Unrecoverable error: %s (%s)", u.Status, u.Message)
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
