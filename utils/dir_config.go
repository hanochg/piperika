package utils

import "fmt"

var constDirConfig = map[string]*DirConfig{ // TODO use .rc file like in each directory instead of const / or use pipeline yaml
	"/access": {
		PipelineName:      "access_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/access/server": {
		PipelineName:      "access_build",
		DefaultStep:       "access_server",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/access/client": {
		PipelineName:      "access_build",
		DefaultStep:       "access_client",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/access-nodejs-client": {
		PipelineName:      "access_build",
		DefaultStep:       "access_nodejs_client",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/access-go-client": {
		PipelineName:      "access_build",
		DefaultStep:       "access_go_client",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
}

type DirConfig struct {
	PipelineName      string
	DefaultStep       string
	PipelinesSourceId int
}

func GetDirConfig() (*DirConfig, error) {
	dir, err := GetRelativeDir()
	if err != nil {
		return nil, err
	}
	if config, ok := constDirConfig[dir]; !ok {
		return nil, fmt.Errorf("working directory '%s' is not mapped, could not resolve pipeline to work against", dir)
	} else {
		return config, nil
	}
}
