package commands

import (
	"context"
	"github.com/hanochg/piperika/actions/report"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"time"
)

func PrintReport() components.Command {
	return components.Command{
		Name:        "report",
		Description: "Fetch multiple Pipes data and prints a report",
		Aliases:     []string{"r"},
		Arguments:   getArguments(),
		Flags:       getReportsFlags(),
		Action:      printReport,
	}
}

func getReportsFlags() []components.Flag {
	return []components.Flag{
		plugins.GetServerIdFlag(),
		components.StringFlag{
			Name:         "branch",
			Description:  "The relevant Git branch for the reports",
			DefaultValue: "master",
		},
	}
}

func printReport(c *components.Context) error {
	client, err := http.NewPipelineHttp(c)
	if err != nil {
		return err
	}
	config, err := utils.GetConfigurations()
	if err != nil {
		return err
	}
	uiUrl, err := utils.GetUIBaseUrl(c)
	if err != nil {
		return err
	}
	branch, err := utils.GetCurrentBranchName(c)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.BranchName, branch)
	ctx = context.WithValue(ctx, utils.BaseUiUrl, uiUrl)
	ctx = context.WithValue(ctx, utils.HttpClientCtxKey, client)
	ctx = context.WithValue(ctx, utils.ConfigCtxKey, config)
	return report.ReportsGathering(ctx)
}
