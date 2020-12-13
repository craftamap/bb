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
	Stage  PipelineStateResult
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

type StepImage struct {
	Name string `mapstructure:"name"`
}

type StepCommand struct {
	Name        string `mapstructure:"name"`
	Command     string `mapstructure:"command"`
	Action      string `mapstructure:"action"`
	CommandType string `mapstructure:"commandType"`
}

type Step struct {
	Name              string        `mapstructure:"name"`
	Pipeline          Pipeline      `mapstructure:"pipeline"`
	State             PipelineState `mapstructure:"state"`
	RunNumber         int           `mapstructure:"run_number"`
	CompletedOn       string        `mapstructure:"completed_on"`
	MaxTime           int           `mapstructure:"maxTime"`
	Image             StepImage     `mapstructure:"image"`
	UUID              string        `mapstructure:"uuid"`
	CreatedOn         string        `mapstructure:"created_on"`
	BuildSecondsUsed  int           `mapstructure:"build_seconds_used"`
	DurationInSeconds int           `mapstructure:"duration_in_seconds"`
	TeardownCommands  []StepCommand `mapstructure:"teardown_commands"`
	ScriptCommands    []StepCommand `mapstructure:"script_commands"`
	SetupCommands     []StepCommand `mapstructure:"setup_commands"`
	Type              string        `mapstructure:"type"`
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

func (c Client) PipelineGet(repoOrga string, repoSlug string, idOrUuid string) (*Pipeline, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	response, err := client.Repositories.Pipelines.Get(&bitbucket.PipelinesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		IDOrUuid: idOrUuid,
	})
	if err != nil {
		return nil, err
	}
	var pipeline *Pipeline
	err = mapstructure.Decode(response, &pipeline)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func (c Client) PipelineStepsList(repoOrga string, repoSlug, idOrUuid string) (*[]Step, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	response, err := client.Repositories.Pipelines.ListSteps(&bitbucket.PipelinesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		IDOrUuid: idOrUuid,
	})
	if err != nil {
		return nil, err
	}
	var steps *[]Step
	err = mapstructure.Decode(response.(map[string]interface{})["values"], &steps)
	if err != nil {
		return nil, err
	}
	return steps, nil
}

func (c Client) PipelinesLogs(repoOrga string, repoSlug, idOrUuid string, StepUuid string) (string, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	return client.Repositories.Pipelines.GetLog(&bitbucket.PipelinesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		IDOrUuid: idOrUuid,
		StepUuid: StepUuid,
	})

}
