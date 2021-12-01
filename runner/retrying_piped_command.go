package runner

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"time"
)

type runner interface {
	init(ctx context.Context, state *pipedCommandState) (string, error)
	tick(ctx context.Context, state *pipedCommandState) (*tickStatus, error)
	timedOutOperation(ctx context.Context, state *pipedCommandState, status *tickStatus) (string, error)
}

type runnerConfig struct {
	interval time.Duration
	timeout  time.Duration
}

func newRetryingPipedCommand(operationName string, runner runner, runnerConfig runnerConfig) *waitedPipedCommand {
	return &waitedPipedCommand{
		runner:        runner,
		runnerConfig:  runnerConfig,
		operationName: operationName,
	}
}

type waitedPipedCommand struct {
	runner
	operationName string
	runnerConfig
}

func (c *waitedPipedCommand) Run(ctx context.Context, state *pipedCommandState) error {
	maxRetries := c.timeout.Milliseconds() / c.interval.Milliseconds()
	backoffConfig := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(maxRetries))
	failBackoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(c.interval), uint64(maxRetries))

	lastStatus, err := c.ticker(ctx, state, backoffConfig, failBackoff)
	if err != nil {
		if !errors.As(err, &timeOutError{}) {
			return err
		}
		fmt.Printf("%s: timed-out\n", c.operationName)
	}

	message, err := c.timedOutOperation(ctx, state, lastStatus)
	if err != nil {
		return err
	}

	fmt.Printf("%s: %s\n", c.operationName, message)

	_, err = c.ticker(ctx, state, backoffConfig, failBackoff)
	return err
}

func (c *waitedPipedCommand) ticker(ctx context.Context, state *pipedCommandState, backoffConfig backoff.BackOff, failBackoff backoff.BackOff) (*tickStatus, error) {
	var currentStatus *tickStatus
	for currInterval := backoffConfig.NextBackOff(); currInterval > 0; currInterval = backoffConfig.NextBackOff() {
		err := backoff.RetryNotify(func() (err error) {
			currentStatus, err = c.tick(ctx, state)
			return
		}, failBackoff, func(err error, duration time.Duration) {
			// TODO error handling
		})
		if err != nil {
			return currentStatus, err
		}

		// TODO Handle terminal using status
		fmt.Printf("%s: %s - %s\n", c.operationName, currentStatus.status, currentStatus.message)

		if currentStatus.done {
			return currentStatus, nil
		}
	}
	return currentStatus, fmt.Errorf("timed-out %w", timeOutError{})
}

type timeOutError struct {
	error
}

type tickStatus struct {
	status  string
	message string
	done    bool
}
