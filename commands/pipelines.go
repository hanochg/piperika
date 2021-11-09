package commands

import (
	"github.com/jfrog/jfrog-cli-core/plugins/components"
)

func GetCommand() components.Command {
	return components.Command{
		Name:        "pipelines",
		Description: "pipelines operations",
		Aliases:     []string{"pl"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		Action:      action,
	}
}

func getArguments() []components.Argument {
	return []components.Argument{}
}

func getFlags() []components.Flag {
	return []components.Flag{}
}

func action(c *components.Context) error {
	return nil
}
