package runner

import "time"

var mainDefaultConfig = runnerConfig{interval: time.Second, timeout: time.Minute * 10}

var registry = []pipedCommand{
	newRetryingPipedCommand("fetch branch", new001GetPipeSourceBranch(), mainDefaultConfig),
	newRetryingPipedCommand("sync", new002WaitPipSourceCompletion(), mainDefaultConfig),
	newRetryingPipedCommand("wait for run creation", new003GetRun(), mainDefaultConfig),
	newRetryingPipedCommand("follow run", new004WaitForRun(), mainDefaultConfig),
}
