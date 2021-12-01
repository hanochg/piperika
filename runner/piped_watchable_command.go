package runner

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/hanochg/piperika/runner/datastruct"
	"time"
)

type runnerConfig struct {
	interval time.Duration
	timeout  time.Duration
}

type watchablePipedCommand struct {
	datastruct.Runner
	operationName string
	runnerConfig
}

type timeOutError struct {
	error
}

func NewWatchablePipedCommand(operationName string, runner datastruct.Runner, runnerConfig runnerConfig) *watchablePipedCommand {
	return &watchablePipedCommand{
		Runner:        runner,
		runnerConfig:  runnerConfig,
		operationName: operationName,
	}
}

func (c *watchablePipedCommand) Run(ctx context.Context, state *datastruct.PipedCommandState) error {
	maxRetries := c.timeout.Milliseconds() / c.interval.Milliseconds()
	backoffConfig := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(maxRetries))
	failBackoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(maxRetries))

	_, err := c.Init(ctx, state)
	if err != nil {
		return err
	}

	lastStatus, err := c.ticker(ctx, state, backoffConfig, failBackoff)
	if err != nil {
		if errors.As(err, &timeOutError{}) {
			fmt.Printf("%s: timed-out\n", c.operationName)
		}
		return err
	}

	message, err := c.OnComplete(ctx, state, lastStatus)
	if err != nil {
		return err
	}
	fmt.Printf("%s: %s\n", c.operationName, message)
	return nil
}

func (c *watchablePipedCommand) ticker(ctx context.Context, state *datastruct.PipedCommandState, backoffConfig backoff.BackOff, failBackoff backoff.BackOff) (*datastruct.RunStatus, error) {
	var currentStatus *datastruct.RunStatus
	for currInterval := backoffConfig.NextBackOff(); currInterval > 0; currInterval = backoffConfig.NextBackOff() {
		err := backoff.RetryNotify(func() (err error) {
			currentStatus, err = c.Tick(ctx, state)
			return
		}, failBackoff, func(err error, duration time.Duration) {
			// TODO error handling
		})
		if err != nil {
			return currentStatus, err
		}

		// TODO Handle terminal using status
		fmt.Printf("%s: %s - %s\n", c.operationName, currentStatus.Status, currentStatus.Message)

		if currentStatus.Done {
			return currentStatus, nil
		}
	}
	return currentStatus, fmt.Errorf("timed-out %w", timeOutError{})
}
