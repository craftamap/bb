package internal

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/wbrefvem/go-bitbucket"
)

type Repository struct {
	Links       map[string]interface{} `mapstructure:"Links"`
	UUID        string                 `mapstructure:"Uuid"`
	FullName    string                 `mapstructure:"FullName"`
	IsPrivate   bool                   `mapstructure:"IsPrivate"`
	Parent      *Repository            `mapstructure:"Parent"`
	Owner       *Account               `mapstructure:"Owner"`
	Name        string                 `mapstructure:"Name"`
	Description string                 `mapstructure:"Description"`
	CreatedOn   time.Time              `mapstructure:"CreatedOn"`
	UpdatedOn   time.Time              `mapstructure:"UpdatedOn"`
	Size        int                    `mapstructure:"Size"`
	Language    string                 `mapstructure:"Language"`
	HasIssues   bool                   `mapstructure:"HasIssues"`
	HasWiki     bool                   `mapstructure:"HasWiki"`
	ForkPolicy  string                 `mapstructure:"ForkPolicy"`
	MainBranch  *Branch                `mapstructure:"Mainbranch"`
	// Project     Project         `mapstructure:"project"`
}

func (c Client) RepositoryGet(repoOrga string, repoSlug string) (*Repository, error) {
	client := bitbucket.NewAPIClient(bitbucket.NewConfiguration())
	response, _, err := client.RepositoriesApi.RepositoriesUsernameRepoSlugGet(
		context.WithValue(context.Background(), bitbucket.ContextBasicAuth, bitbucket.BasicAuth{
			UserName: c.Username,
			Password: c.Password,
		}), repoOrga, repoSlug)

	if err != nil {
		return nil, err
	}

	var repo Repository
	err = mapstructure.Decode(response, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}
