package runner

import (
	"context"
)

func RunPipe(ctx context.Context) error {
	pipedState := &pipedCommandState{}
	for _, cmd := range registry {
		err := cmd.Run(ctx, pipedState)
		if err != nil {
			return err
		}
	}
	return nil
}
