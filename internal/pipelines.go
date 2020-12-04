package internal

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type PipelineStateResult struct {
	Name string
	Type string
}

type PipelineState struct {
	Name   string
	Type   string
	Result PipelineStateResult
}

type PipelineTrigger struct {
	Name string
	Type string
}

type Pipeline struct {
	Type              string
	UUID              string        `mapstructure:"uuid"`
	PipelineState     PipelineState `mapstructure:"state"`
	BuildNumber       int           `mapstructure:"build_number"`
	Creator           Account
	CreatedOn         string `mapstructure:"created_on"`
	CompletedOn       string `mapstructure:"completed_on"`
	Target            interface{}
	Trigger           PipelineTrigger
	RunNumber         int  `mapstructure:"run_number"`
	DurationInSeconds int  `mapstructure:"duration_in_seconds"`
	BuildSecondsUsed  int  `mapstructure:"build_seconds_used"`
	FirstSuccessful   bool `mapstructure:"first_successful"`
	Expired           bool
}

func (c Client) PipelineList(repoOrga string, repoSlug string) (*[]Pipeline, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	response, err := client.Repositories.Pipelines.List(&bitbucket.PipelinesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		Sort:     "-created_on",
	})
	if err != nil {
		return nil, err
	}
	var pipelines *[]Pipeline
	err = mapstructure.Decode(response.(map[string]interface{})["values"], &pipelines)
	if err != nil {
		return nil, err
	}
	return pipelines, nil
}
