package internal

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type Commit struct {
	Hash       string                 `mapstructure:"hash"`
	Type       string                 `mapstructure:"type"`
	Message    string                 `mapstructure:"message"`
	Parents    []*Commit              `mapstructure:"parents"`
	Repository *Repository            `mapstructure:"repository"`
	Author     map[string]interface{} `mapstructure:"author"`
}

type Commits struct {
	Values []*Commit `mapstructure:"values"`
}

func (c Client) GetCommits(repoOrga string, repoSlug string, revision string, include string, exclude string) (*Commits, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opts := bitbucket.CommitsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		Revision: revision,
		Exclude:  exclude,
		Include:  include,
	}

	var commits Commits
	response, err := client.Repositories.Commits.GetCommits(&opts)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(response, &commits)
	if err != nil {
		return nil, err
	}
	return &commits, nil
}
