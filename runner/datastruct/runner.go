package datastruct

import "context"

type Runner interface {
	Init(ctx context.Context, state *PipedCommandState) (string, error)
	Tick(ctx context.Context, state *PipedCommandState) (*RunStatus, error)
	OnComplete(ctx context.Context, state *PipedCommandState, status *RunStatus) (string, error)
}
