package runner

import (
	"context"
	"github.com/hanochg/piperika/runner/command"
)

func RunPipe(ctx context.Context, pipe ...command.PipedCommand) error {
	pipedState := &command.PipedCommandState{}
	for _, cmd := range pipe {
		res := cmd.Run(ctx, pipedState)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}
