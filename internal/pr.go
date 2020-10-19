package internal

import (
	"fmt"

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
	ID                int      `mapstructure:"id"`
	Title             string   `mapstructure:"title"`
	State             string   `mapstructure:"state"`
	Source            Resource `mapstructure:"source"`
	Destination       Resource `mapstructure:"destination"`
	Type              string   `mapstructure:"type"`
	TaskCount         int      `mapstructure:"task_count"`
	Description       string   `mapstructure:"description"`
	Author            User     `mapstructure:"author"`
	CloseSourceBranch bool     `mapstructure:"close_source_branch"`
	CommentCount      int      `mapstructure:"comment_count"`
	CreatedOn         string   `mapstructure:"created_on"`
	MergeCommit       Commit   `mapstructure:"merge_commit"`
}

type Resource struct {
	Branch     Branch     `mapstructure:"branch"`
	Commit     Commit     `mapstructure:"commit"`
	Repository Repository `mapstructure:"repository"`
}

type Branch struct {
	Name string `mapstructure:"name"`
}

type Commit struct {
	Hash string `mapstructure:"hash"`
	Type string `mapstructure:"type"`
}

type Repository struct {
	FullName string `mapstructure:"full_name"`
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`
	UUID     string `mapstructure:"uuid"`
}

type User struct {
	AccountID   string `mapstructure:"account_id"`
	DisplayName string `mapstructure:"display_name"`
	Nickname    string `mapstructure:"nickname"`
	Type        string `mapstructure:"user"`
	UUID        string `mapstructure:"uuid"`
}

func PrList(username string, password string, repoOrga string, repoSlug string) (*ListPullRequests, error) {
	client := bitbucket.NewBasicAuth(username, password)

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

func GetPrIDBySourceBranch(username string, password string, repoOrga string, repoSlug string, sourceBranch string) (*ListPullRequests, error) {
	client := bitbucket.NewBasicAuth(username, password)

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

func PrView(username string, password string, repoOrga string, repoSlug string, id string) (*PullRequest, error) {
	client := bitbucket.NewBasicAuth(username, password)

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
