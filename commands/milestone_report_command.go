package commands

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/report"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"sync"
	"time"
)

const defaultBranchName string = "master"

var services = []string{"access", "metadata", "event", "router", "jfconnect", "mission_control", "integration"}

func GetMilestoneReport() components.Command {
	return components.Command{
		Name:        "milestone_report",
		Description: "Get milestone report for Infra QA blander",
		Aliases:     []string{"m"},
		Arguments:   getArguments(),
		Flags:       getFlags(),
		Action:      getMilestoneReport,
	}
}

func getMilestoneReport(c *components.Context) error {
	uiUrl, err := utils.GetUIBaseUrl(c)
	if err != nil {
		return err
	}
	httpClient, err2 := getHttpClient(c)
	if err2 != nil {
		return err2
	}

	var wg sync.WaitGroup
	wg.Add(len(services))

	for i := 0; i < len(services); i++ {
		go func(i int) {
			defer wg.Done()
			serviceReport := report.GetServiceReport(httpClient, services[i], uiUrl, defaultBranchName)
			fmt.Println(serviceReport.String())
		}(i)
	}
	wg.Wait()
	fmt.Println("Finished")
	return nil
}

func getHttpClient(c *components.Context) (http.PipelineHttpClient, error) {
	client, err := http.NewPipelineHttp(c)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.HttpClientCtxKey, client)
	httpClient := ctx.Value(utils.HttpClientCtxKey).(http.PipelineHttpClient)
	return httpClient, err
}
