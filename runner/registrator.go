package runner

import (
	"github.com/hanochg/piperika/runner/datastruct"
	"github.com/hanochg/piperika/runner/runs"
	"time"
)

var mainDefaultConfig = runnerConfig{interval: time.Second, timeout: time.Minute * 10}

var registry = []datastruct.PipedCommand{
	NewPipedCommand("fetch branch", runs.New001GetPipeSourceBranch(), mainDefaultConfig),
	NewWatchablePipedCommand("sync", runs.New002WaitPipSourceCompletion(), mainDefaultConfig),
	NewPipedCommand("create or grab current run", runs.New003GetRun(), mainDefaultConfig),
	NewWatchablePipedCommand("follow run", runs.New004WaitForRun(), mainDefaultConfig),
}
