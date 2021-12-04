package utils

import "fmt"

var constDirConfig = map[string]*DirConfig{ // TODO use .rc file like in each directory instead of const
	"/access": {
		PipelineName:      "access_build",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/access-go-client": {
		PipelineName:      "access_build",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
}

type DirConfig struct {
	PipelineName      string
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
