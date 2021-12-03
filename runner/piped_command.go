package runner

import (
	"context"
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

	if waitErr != nil {
		triggerErr := c.TriggerStateChange(ctx, state)
		if triggerErr != nil {
			return triggerErr
		}

		secondWaitErr := c.retryResolveState(ctx, state)
		if secondWaitErr != nil {
			return secondWaitErr
		}
	}

	fmt.Printf("%s: finished\n", c.operationName)
	return nil
}

func (c *retryingPipedCommand) retryResolveState(ctx context.Context, state *command.PipedCommandState) error {
	return backoff.RetryNotify(
		func() error {
			status, err := c.ResolveState(ctx, state)
			if err != nil {
				return err
			}
			return terminal.UpdateStatus(c.operationName, status.PipelinesStatus, status.Message, "TBD")
		},
		c.newBackoffContext(ctx),
		func(err error, duration time.Duration) {
			println(fmt.Sprintf("######### %v", err.Error()))
		},
	)
}

func (c *retryingPipedCommand) newBackoffContext(ctx context.Context) backoff.BackOffContext {
	initialBackoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(c.maxRetries))
	return backoff.WithContext(initialBackoff, ctx)
}
