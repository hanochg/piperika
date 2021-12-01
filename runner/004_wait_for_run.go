package runner

import (
	"context"
)

func (_ _04) init(ctx context.Context, state *pipedCommandState) (string, error) {
	return "", nil
}

func (_ _04) tick(ctx context.Context, state *pipedCommandState) (*tickStatus, error) {
	return nil, nil
}

func (_ _04) timedOutOperation(ctx context.Context, state *pipedCommandState, status *tickStatus) (string, error) {
	return "", nil
}

type _04 struct {
}

func new004WaitForRun() _04 {
	return _04{}
}
