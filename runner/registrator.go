package runner

import "time"

var mainDefaultConfig = runnerConfig{interval: time.Second, timeout: time.Minute * 10}

var registry = []pipedCommand{
	newRetryingPipedCommand("sync", new001GetPipeSourceBranch(), mainDefaultConfig),
}
