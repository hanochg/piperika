package commands

import (
	"context"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/runner"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"time"
)

func theCommand(c *components.Context) error { // TODO think of better name
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()
	ctx = context.WithValue(ctx, "client", client)
	return runner.RunPipe(ctx)
}
