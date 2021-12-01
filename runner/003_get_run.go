package runner

import (
	"context"
)

func (_ _03) init(ctx context.Context, state *pipedCommandState) (string, error) {
	return "", nil
}

func (_ _03) tick(ctx context.Context, state *pipedCommandState) (*tickStatus, error) {
	return nil, nil
}

func (_ _03) timedOutOperation(ctx context.Context, state *pipedCommandState, status *tickStatus) (string, error) {
	return "", nil
}

type _03 struct {
}

func new003GetRun() _03 {
	return _03{}
}
