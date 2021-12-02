package runner

import (
	"context"
	"github.com/hanochg/piperika/runner/datastruct"
)

func RunPipe(ctx context.Context) error {
	pipedState := &datastruct.PipedCommandState{}
	for _, cmd := range registry {
		err := cmd.Run(ctx, pipedState)
		if err != nil {
			return err
		}
	}
	return nil
}
