package client

import (
	"strings"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type ListIssues struct {
	Values []Issue `mapstructure:"values"`
	Size   int     `mapstructure:"size"`
}

type Issue struct {
	ID         int             `mapstructure:"id"`
	Priority   string          `mapstructure:"priority"`
	Kind       string          `mapstructure:"kind"`
	Type       string          `mapstructure:"type"`
	Repository Repository      `mapstructure:"repository"`
	Links      map[string]Link `mapstructure:"links"`
	Reporter   Account         `mapstructure:"reporter"`
	Title      string          `mapstructure:"title"`
	Component  Component       `mapstructure:"component"`
	Votes      int64           `mapstructure:"votes"`
	Watches    int64           `mapstructure:"watches"`
	Content    CommentContent  `mapstructure:"content"`
	Assignee   Account         `mapstructure:"assignee"`
	State      string          `mapstructure:"state"`
	Version    Version         `mapstructure:"version"`
	CreatedOn  string          `mapstructure:"created_on"`
	EditedOn   string          `mapstructure:"edited_on"`
	UpdatedOn  string          `mapstructure:"edited_on"`
	Milestone  Milestone       `mapstructure:"milestone"`
}

type Version struct {
	Name  string          `mapstructure:"name"`
	Links map[string]Link `mapstructure:"links"`
}

type Milestone struct {
	Name  string          `mapstructure:"name"`
	Links map[string]Link `mapstructure:"links"`
}

type Component struct {
	Name  string          `mapstructure:"name"`
	Links map[string]Link `mapstructure:"links"`
}

type IssueComments struct {
	Values []IssueComment `mapstructure:"values"`
	Size   int            `mapstructure:"size"`
}

type IssueComment struct {
	ID        int             `mapstructure:"id"`
	Type      string          `mapstructure:"type"`
	Links     map[string]Link `mapstructure:"links"`
	Issue     Issue           `mapstructure:"issue"`
	Content   CommentContent  `mapstructure:"content"`
	CreatedOn string          `mapstructure:"created_on"`
	User      Account         `mapstructure:"user"`
	UpdatedOn string          `mapstructure:"edited_on"`
}

type IssueChanges struct {
	Values []IssueChange `mapstructure:"values"`
	Size   int           `mapstructure:"size"`
}
type IssueChange struct {
	ID        int               `mapstructure:"id"`
	Type      string            `mapstructure:"type"`
	Links     map[string]Link   `mapstructure:"links"`
	Issue     Issue             `mapstructure:"issue"`
	CreatedOn string            `mapstructure:"created_on"`
	User      Account           `mapstructure:"user"`
	Changes   map[string]Change `mapstructure:"changes"`
	Message   CommentContent    `mapstructure:"message"`
}

type Change struct {
	New string `mapstructure:"new"`
	Old string `mapstructure:"old"`
}

func (c Client) IssuesList(repoOrga string, repoSlug string, states []string) (*ListIssues, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)
	var query strings.Builder
	query.WriteString("(")
	for i, state := range states {
		query.WriteString("state = \"")
		query.WriteString(state)
		query.WriteString("\" ")
		if i != len(states)-1 {
			query.WriteString("OR ")
		}
	}
	query.WriteString(")")

	opts := bitbucket.IssuesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		Query:    query.String(),
	}

	response, err := client.Repositories.Issues.Gets(&opts)
	if err != nil {
		return nil, err
	}

	var issues ListIssues
	err = mapstructure.Decode(response, &issues)
	if err != nil {
		return nil, err
	}

	return &issues, nil
}

func (c Client) IssuesView(repoOrga string, repoSlug string, id string) (*Issue, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opts := bitbucket.IssuesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.Issues.Get(&opts)
	if err != nil {
		return nil, err
	}

	var issue Issue
	err = mapstructure.Decode(response, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c Client) IssuesViewComments(repoOrga string, repoSlug string, id string) (*IssueComments, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opts := bitbucket.IssuesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.Issues.GetComments(&bitbucket.IssueCommentsOptions{
		IssuesOptions: opts,
	})
	if err != nil {
		return nil, err
	}

	var issueComments IssueComments
	err = mapstructure.Decode(response, &issueComments)
	if err != nil {
		return nil, err
	}

	return &issueComments, nil
}

func (c Client) IssuesViewChanges(repoOrga string, repoSlug string, id string) (*IssueChanges, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)
	opts := bitbucket.IssuesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.Issues.GetChanges(&bitbucket.IssueChangesOptions{
		IssuesOptions: opts,
	})
	if err != nil {
		return nil, err
	}

	var issueChanges IssueChanges
	err = mapstructure.Decode(response, &issueChanges)
	if err != nil {
		return nil, err
	}

	return &issueChanges, nil
}

func (c Client) IssuesCreate(repoOrga string, repoSlug string, options struct {
	Title       string
	Description string
	Assignee    string
	Kind        string
	Priority    string
	Status      string
	Version     string
	Milestone   string
	Component   string
}) (*Issue, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)
	opts := bitbucket.IssuesOptions{
		Owner:     repoOrga,
		RepoSlug:  repoSlug,
		Title:     options.Title,
		Content:   options.Description,
		Assignee:  options.Assignee,
		Kind:      options.Kind,
		Priority:  options.Priority,
		State:     options.Status,
		Version:   options.Version,
		Milestone: options.Milestone,
		Component: options.Component,
	}

	response, err := client.Repositories.Issues.Create(&opts)
	if err != nil {
		return nil, err
	}

	var issue Issue
	err = mapstructure.Decode(response, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c Client) IssuesEdit(repoOrga string, repoSlug string, id string, options struct {
	Title       string
	Description string
	Assignee    string
	Kind        string
	Priority    string
	Status      string
	Version     string
	Milestone   string
	Component   string
}) (*Issue, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opts := bitbucket.IssuesOptions{
		ID:        id,
		Owner:     repoOrga,
		RepoSlug:  repoSlug,
		Title:     options.Title,
		Content:   options.Description,
		Assignee:  options.Assignee,
		Kind:      options.Kind,
		Priority:  options.Priority,
		State:     options.Status,
		Version:   options.Version,
		Milestone: options.Milestone,
		Component: options.Component,
	}

	response, err := client.Repositories.Issues.Update(&opts)
	if err != nil {
		return nil, err
	}

	var issue Issue
	err = mapstructure.Decode(response, &issue)
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c Client) IssuesComment(repoOrga string, repoSlug string, id string, message string) (*IssueComment, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opts := bitbucket.IssuesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.Issues.CreateComment(&bitbucket.IssueCommentsOptions{
		IssuesOptions:  opts,
		CommentContent: message,
	})
	if err != nil {
		return nil, err
	}

	var issueComment IssueComment
	err = mapstructure.Decode(response, &issueComment)
	return &issueComment, err
}

func (c Client) IssuesDelete(repoOrga string, repoSlug string, id string) (interface{}, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opts := bitbucket.IssuesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	return client.Repositories.Issues.Delete(&opts)
}
