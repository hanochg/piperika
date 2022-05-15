package requests

import (
	"encoding/json"
	"github.com/hanochg/piperika/http"
)

const (
	projectsUrl = "/projects"
)

type GetProjectsOptions struct {
	ProjectIds string `url:"projectIds,omitempty"`
	Names      string `url:"names,omitempty"`
}

type Project struct {
	Name       string `json:"name"`
	ProjectIds int    `json:"id"`
	SourceId   string `json:"sourceId"`
}

type ProjectsResponse struct {
	Projects []Project
}

func GetProjects(client http.PipelineHttpClient, options GetProjectsOptions) (*ProjectsResponse, error) {
	body, err := client.SendGet(projectsUrl, http.ClientOptions{Query: options})
	if err != nil {
		return nil, err
	}
	res := &ProjectsResponse{}
	err = json.Unmarshal(body, &res.Projects)
	return res, err
}
