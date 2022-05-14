package requests

import (
	"encoding/json"
	"fmt"
	"github.com/hanochg/piperika/http"
)

const (
	projectsUrl = "/projects"
)

type ProjectsOptions struct {
	ProjectId int
}

type ProjectsResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetProjects(client http.PipelineHttpClient, options ProjectsOptions) (*ProjectsResponse, error) {
	body, err := client.SendGet(projectsUrl+fmt.Sprintf("/%d", options.ProjectId), http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &ProjectsResponse{}
	err = json.Unmarshal(body, &res)
	return res, err
}
