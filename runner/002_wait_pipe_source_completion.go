package runner

import (
	"context"
)

func (_ _02) init(ctx context.Context, state *pipedCommandState) (string, error) {
	return "", nil
}

func (_ _02) tick(ctx context.Context, state *pipedCommandState) (*tickStatus, error) {
	return nil, nil
}

func (_ _02) timedOutOperation(ctx context.Context, state *pipedCommandState, status *tickStatus) (string, error) {
	return "", nil
}

type _02 struct {
}

func new002WaitPipSourceCompletion() _02 {
	return _02{}
}
