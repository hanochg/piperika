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
	"/go-base-application": {
		PipelineName:      "go_base_application_build",
		DefaultStep:       "go_base_application",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},

	"/artifactory": {
		PipelineName:      "artifactory_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/jfconnect": {
		PipelineName:      "jfconnect_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/jfconnect/service": {
		PipelineName:      "jfconnect_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/jfconnect/go-client": {
		PipelineName:      "jfconnect_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/jfconnect/java-client": {
		PipelineName:      "jfconnect_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/jfconnect/nodejs-client": {
		PipelineName:      "jfconnect_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/metadata": {
		PipelineName:      "metadata_build",
		DefaultStep:       "trigger_all",
		PipelinesSourceId: ArtifactoryPipelinesSourceId,
	},
	"/router": {
		PipelineName:      "router_build",
		DefaultStep:       "trigger_all",
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
