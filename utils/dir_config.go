package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var PIPERIKA_CONFIG_FILE = ".pipr"

type Configurations struct {
	PipelineName      string `json:"pipeline_name,omitempty"`
	DefaultStep       string `json:"default_step,omitempty"`
	PipelinesSourceId int    `json:"pipelines_source_id,omitempty"`
}

func GetConfigurations() (*Configurations, error) {
	gitRootDir, err := GetRootDir()
	if err != nil {
		return nil, err
	}

	curDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	defaultConfig := &Configurations{
		PipelineName:      "",
		DefaultStep:       "",
		PipelinesSourceId: 0,
	}
	conf, err := loadConfigFromFolder(curDir, defaultConfig)
	for err != nil && curDir != gitRootDir {
		curDir = filepath.Dir(curDir)
		conf, err = loadConfigFromFolder(curDir, defaultConfig)
	}

	if err != nil {
		return nil, err
	}

	return conf, nil
}

func loadConfigFromFolder(dir string, defaultConfig *Configurations) (*Configurations, error) {
	piprConfFile := filepath.Join(dir, PIPERIKA_CONFIG_FILE)
	if _, err := os.Stat(piprConfFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("cannot find '.pipr' configuration file")
	}
	confFile, err := os.Open(piprConfFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open configuration file, path [%v] err [%v]", piprConfFile, err)
	}
	defer func() { _ = confFile.Close() }()

	confFileBytes, err := ioutil.ReadAll(confFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read configuration file content, path [%v] err [%v]", piprConfFile, err)
	}

	if err := json.Unmarshal(confFileBytes, &defaultConfig); err != nil {
		return nil, fmt.Errorf("cannot parse configuration file, path [%v] err [%v]", piprConfFile, err)
	}
	return defaultConfig, nil
}
