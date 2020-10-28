package internal

import (
	"fmt"
	"strings"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type ListPullRequests struct {
	Size     int           `mapstructure:"size"`
	Page     int           `mapstructure:"page"`
	PageLen  int           `mapstructure:"pagelen"`
	Next     string        `mapstructure:"next"`
	Previous string        `mapstructure:"previous"`
	Values   []PullRequest `mapstructure:"values"`
}

type PullRequest struct {
	ID                int             `mapstructure:"id"`
	Title             string          `mapstructure:"title"`
	State             string          `mapstructure:"state"`
	Source            Resource        `mapstructure:"source"`
	Destination       Resource        `mapstructure:"destination"`
	Type              string          `mapstructure:"type"`
	TaskCount         int             `mapstructure:"task_count"`
	Description       string          `mapstructure:"description"`
	Author            Account         `mapstructure:"author"`
	CloseSourceBranch bool            `mapstructure:"close_source_branch"`
	CommentCount      int             `mapstructure:"comment_count"`
	CreatedOn         string          `mapstructure:"created_on"`
	MergeCommit       Commit          `mapstructure:"merge_commit"`
	Links             map[string]Link `mapstructure:"links"`
}

type Resource struct {
	Branch     Branch     `mapstructure:"branch"`
	Commit     Commit     `mapstructure:"commit"`
	Repository Repository `mapstructure:"repository"`
}

type Status struct {
	Type        string                 `mapstructure:"type"`
	Links       map[string]interface{} `mapstructure:"links"`
	UUID        string                 `mapstructure:"uuid"`
	Key         string                 `mapstructure:"key"`
	Refname     string                 `mapstructure:"refname"`
	URL         string                 `mapstructure:"url"`
	State       string                 `mapstructure:"state"`
	Name        string                 `mapstructure:"name"`
	Description string                 `mapstructure:"description"`
	CreatedOn   string                 `mapstructure:"created_on"`
	UpdatedOn   string                 `mapstructure:"updated_on"`
}

type Statuses struct {
	Size     int      `mapstructure:"size"`
	Page     int      `mapstructure:"page"`
	PageLen  int      `mapstructure:"pagelen"`
	Next     string   `mapstructure:"next"`
	Previous string   `mapstructure:"previous"`
	Values   []Status `mapstructure:"values"`
}

func (c Client) PrList(repoOrga string, repoSlug string) (*ListPullRequests, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
	}

	response, err := client.Repositories.PullRequests.Gets(opt)
	if err != nil {
		return nil, err
	}

	var pullRequests ListPullRequests
	err = mapstructure.Decode(response, &pullRequests)
	if err != nil {
		return nil, err
	}

	return &pullRequests, nil
}

func (c Client) GetPrIDBySourceBranch(repoOrga string, repoSlug string, sourceBranch string) (*ListPullRequests, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		Query:    fmt.Sprintf("source.branch.name = \"%s\"", sourceBranch),
	}

	response, err := client.Repositories.PullRequests.Gets(opt)
	if err != nil {
		return nil, err
	}

	var pullRequests ListPullRequests
	err = mapstructure.Decode(response, &pullRequests)
	if err != nil {
		return nil, err
	}

	return &pullRequests, nil
}

func (c Client) PrView(repoOrga string, repoSlug string, id string) (*PullRequest, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.PullRequests.Get(opt)
	if err != nil {
		return nil, err
	}

	var pullRequest PullRequest
	err = mapstructure.Decode(response, &pullRequest)
	if err != nil {
		return nil, err
	}
	return &pullRequest, nil
}

func (c Client) PrCreate(repoOrga string, repoSlug string, sourceBranch string, destinationBranch string, title string, body string, reviewers []string) (*PullRequest, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:             repoOrga,
		RepoSlug:          repoSlug,
		SourceBranch:      sourceBranch,
		DestinationBranch: destinationBranch,
		Title:             title,
		Description:       body,
		Reviewers:         reviewers,
	}

	response, err := client.Repositories.PullRequests.Create(opt)

	if err != nil {
		return nil, err
	}

	var pullRequest PullRequest
	err = mapstructure.Decode(response, &pullRequest)
	if err != nil {
		return nil, err
	}
	return &pullRequest, nil
}

func (c Client) PrStatuses(repoOrga string, repoSlug string, id string) (*Statuses, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)
	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.PullRequests.Statuses(opt)
	if err != nil {
		return nil, err
	}

	var statuses Statuses
	err = mapstructure.Decode(response, &statuses)
	if err != nil {
		return nil, err
	}

	return &statuses, nil

}

func (c Client) PrDefaultBody(repoOrga string, repoSlug string, sourceBranch string, destinationBranch string) (string, error) {
	commits, err := c.GetCommits(repoOrga, repoSlug, sourceBranch, "", destinationBranch)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, commit := range commits.Values {
		sb.WriteString("- " + commit.Message + "\n")
	}

	return sb.String(), nil
}
