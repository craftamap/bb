package internal

import (
	"context"
	"encoding/json"
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

type DefaultReviewers struct {
	PageLen int        `json:"pagelen"`
	Values  []*Account `json:"values"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
	Next    string     `json:"next"`
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

func (c Client) GetDefaultReviewers(repoOrga string, repoSlug string) (*DefaultReviewers, error) {
	client := bitbucket.NewAPIClient(bitbucket.NewConfiguration())

	response, err := client.PullrequestsApi.RepositoriesUsernameRepoSlugDefaultReviewersGet(
		context.WithValue(context.Background(), bitbucket.ContextBasicAuth, bitbucket.BasicAuth{
			UserName: c.Username,
			Password: c.Password,
		}), repoOrga, repoSlug)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	defaultReviewers := DefaultReviewers{}
	json.NewDecoder(response.Body).Decode(&defaultReviewers)

	return &defaultReviewers, nil
}
