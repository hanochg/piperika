package runs

import (
	"context"
	"github.com/hanochg/piperika/runner/datastruct"
)

func (_ _04) Init(ctx context.Context, state *datastruct.PipedCommandState) (string, error) {
	return "", nil
}

func (_ _04) Tick(ctx context.Context, state *datastruct.PipedCommandState) (*datastruct.RunStatus, error) {
	return nil, nil
}

func (_ _04) OnComplete(ctx context.Context, state *datastruct.PipedCommandState, status *datastruct.RunStatus) (string, error) {
	return "", nil
}

type _04 struct {
}

func New004WaitForRun() datastruct.Runner {
	return _04{}
}
