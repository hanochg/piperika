package commands

import (
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
)

func GetCommand() components.Command {
	return components.Command{
		Name:        "pipelines",
		Description: "pipelines operations",
		Aliases:     []string{"pp"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		Action:      action,
	}
}

func getArguments() []components.Argument {
	return []components.Argument{
		{Name: "command", Description: "perform a command like get / list"},
	}
}

func getFlags() []components.Flag {
	return []components.Flag{plugins.GetServerIdFlag()}
}

func action(c *components.Context) error {
	return theCommand(c)
}
