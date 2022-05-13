package commands

import (
	"context"
	"fmt"
	"github.com/hanochg/piperika/http"
	"github.com/hanochg/piperika/report"
	"github.com/hanochg/piperika/utils"
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"sync"
	"time"
)

func PrintReport() components.Command {
	return components.Command{
		Name:        "report",
		Description: "Prints a custom user report",
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
			Name:        "name",
			Description: "infra_milestone - Infra services Milestone report. Returns status of latest run for release, build and post-release pipelines.",
			Mandatory:   true,
		},
		components.StringFlag{
			Name:         "branch",
			Description:  "Release pipeline branch",
			DefaultValue: "master",
		},
	}
}

func printReport(c *components.Context) error {
	reportToPrint := selectReportToPrint(c)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for i := 0; i < len(reportToPrint); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			report.FetchReport(reportToPrint[i])
			mutex.Lock()
			fmt.Println(report.ToString(reportToPrint[i]))
			mutex.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Println("Finished")
	return nil
}

func selectReportToPrint(c *components.Context) []report.Report {
	baseUrl := getBaseUrl(c)
	httpClient := getHttpClient(c)
	branch := c.GetStringFlagValue("branch")

	if c.GetStringFlagValue("name") == "infra_milestone" {
		return infraServicesMilestoneReport(httpClient, baseUrl, branch)
	}
	panic("Unknown report name: " + c.GetStringFlagValue("name") + ". Allowed names: 'infra_milestone'")
}

func infraServicesMilestoneReport(httpClient http.PipelineHttpClient, baseUrl string, branch string) []report.Report {
	var rep []report.Report
	for i := 0; i < len(utils.InfraReportServices); i++ {
		el := &report.ServiceReport{
			HttpClient:      httpClient,
			BaseUrl:         baseUrl,
			MilestoneBranch: branch,
			ServiceName:     utils.InfraReportServices[i],
		}
		rep = append(rep, el)
	}
	return rep
}

func getBaseUrl(c *components.Context) string {
	baseUrl, err := utils.GetUIBaseUrl(c)
	if err != nil {
		panic("Unable to get base URL!")
	}
	return baseUrl
}

func getHttpClient(c *components.Context) http.PipelineHttpClient {
	httpClient, err := http.NewPipelineHttp(c)
	if err != nil {
		panic("Unable to get HTTP client!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.HttpClientCtxKey, httpClient)
	return httpClient
}
