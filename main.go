package main

import (
	"github.com/hanochg/piperika/commands"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "piperika"
	app.Description = "JFrog Pipelines in a simple command"
	app.Version = "v1.0.0"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetCommand(),
		commands.GetLinkCommand(),
		commands.PrintReport(),
	}
}
