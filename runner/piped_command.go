package runner

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/runner/datastruct"
)

type pipedCommand struct {
	datastruct.Runner
	operationName string
	runnerConfig
}

func NewPipedCommand(operationName string, runner datastruct.Runner, runnerConfig runnerConfig) *pipedCommand {
	return &pipedCommand{
		Runner:        runner,
		runnerConfig:  runnerConfig,
		operationName: operationName,
	}
}

func (c *pipedCommand) Run(ctx context.Context, state *datastruct.PipedCommandState) error {
	_, err := c.Init(ctx, state)
	if err != nil {
		return err
	}

	_, err = c.Tick(ctx, state)
	if err != nil {
		return err

	}

	fmt.Printf("%s: succeed\n", c.operationName)

	return nil
}
