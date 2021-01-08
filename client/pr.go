package client

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
	Reviewers         []Account       `mapstructure:"reviewers"`
	Participants      []Participant   `mapstructure:"participants"`
	Links             map[string]Link `mapstructure:"links"`
}

type Participant struct {
	Role           string  `mapstructure:"role"`
	State          string  `mapstructure:"state"`
	ParticipatedOn string  `mapstructure:"participated_on"`
	Type           string  `mapstructure:"type"`
	Approved       bool    `mapstructure:"approved"`
	User           Account `mapstructure:"user"`
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

type CommentContent struct {
	Type   string `mapstructure:"type"`
	Raw    string `mapstructure:"raw"`
	Markup string `mapstructure:"markdown"`
	HTML   string `mapstructure:"html"`
}

type CommentParent struct {
	ID      int                    `mapstructure:"id"`
	_Links  map[string]interface{} `mapstructure:"links"`
	Comment *Comment
}

type CommentInline struct {
	To   int    `mapstructure:"to"`
	From int    `mapstructure:"from"`
	Path string `mapstructure:"path"`
}

type Comment struct {
	Links       map[string]interface{} `mapstructure:"links"`
	Parent      CommentParent          `mapstructure:"parent"`
	Deleted     bool                   `mapstructure:"deleted"`
	PullRequest PullRequest            `mapstructure:"pullrequest"`
	Content     CommentContent         `mapstructure:"content"`
	CreatedOn   string                 `mapstructure:"created_on"`
	UpdatedOn   string                 `mapstructure:"updated_on"`
	User        Account                `mapstructure:"user"`
	Inline      CommentInline          `mapstructure:"inline"`
	ID          int                    `mapstructure:"id"`
	Type        string                 `mapstructure:"type"`
	Children    []*Comment
}

type Comments struct {
	Values []*Comment `mapstructure:"values"`
}

func (c Client) PrList(repoOrga string, repoSlug string, states []string) (*ListPullRequests, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		States:   states,
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

func (c Client) PrCreate(
	repoOrga string,
	repoSlug string,
	sourceBranch string,
	destinationBranch string,
	title string,
	body string,
	reviewers []string,
	closeBranch bool) (*PullRequest, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:             repoOrga,
		RepoSlug:          repoSlug,
		SourceBranch:      sourceBranch,
		DestinationBranch: destinationBranch,
		Title:             title,
		Description:       body,
		Reviewers:         reviewers,
		CloseSourceBranch: closeBranch,
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

func (c Client) PrDefaultTitleAndBody(repoOrga string, repoSlug string, sourceBranch string, destinationBranch string) (string, string, error) {
	commits, err := c.GetCommits(repoOrga, repoSlug, sourceBranch, "", destinationBranch)
	if err != nil {
		return "", "", err
	}
	if len(commits.Values) == 0 {
		return sourceBranch, "", nil
	} else if len(commits.Values) == 1 {
		commit := commits.Values[0]

		split := strings.SplitN(commit.Message, "\n", 2)
		if len(split) == 2 {
			return split[0], strings.TrimSpace(split[1]), nil
		} else if len(split) == 1 {
			return split[0], "", nil
		}

		return sourceBranch, "", nil
	} else {
		var sb strings.Builder
		for _, commit := range commits.Values {
			sb.WriteString("- " + strings.Split(commit.Message, "\n")[0] + "\n")
		}

		return sourceBranch, sb.String(), nil
	}
}

func (c Client) PrCommits(repoOrga string, repoSlug string, id string) (*Commits, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}
	response, err := client.Repositories.PullRequests.Commits(opt)
	if err != nil {
		return nil, err
	}

	var commits Commits
	err = mapstructure.Decode(response, &commits)
	if err != nil {
		return nil, err
	}
	return &commits, nil
}

func (c Client) PrMerge(repoOrga string, repoSlug string, id string) (*PullRequest, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.PullRequests.Merge(opt)
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

func (c Client) PrComments(repoOrga string, repoSlug string, id string) (*Comments, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.PullRequests.GetComments(opt)
	if err != nil {
		return nil, err
	}

	var comments Comments
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(response, &comments)
	if err != nil {
		return nil, err
	}
	return &comments, nil
}

func (c Client) PrThreadedComments(repoOrga string, repoSlug string, id string) ([]*Comment, error) {
	comments, err := c.PrComments(repoOrga, repoSlug, id)
	if err != nil {
		return nil, err
	}

	idToComment := map[int]*Comment{}

	for _, comment := range comments.Values {
		idToComment[comment.ID] = comment
	}

	for _, comment := range idToComment {
		parentComment := idToComment[comment.Parent.ID]
		if parentComment != nil {
			// Set parent in child
			comment.Parent.Comment = parentComment
			// Set child in parent
			parentComment.Children = append(parentComment.Children, comment)
		}
	}

	returnArray := []*Comment{}

	for _, comment := range idToComment {
		if comment.Parent.Comment == nil {
			returnArray = append(returnArray, comment)
		}
	}

	return returnArray, nil
}

func (c Client) PrApprove(repoOrga string, repoSlug string, id string) (*Participant, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.PullRequests.Approve(opt)
	if err != nil {
		return nil, err
	}

	var participant Participant
	err = mapstructure.Decode(response, &participant)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (c Client) PrRequestChanges(repoOrga string, repoSlug string, id string) (*Participant, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	response, err := client.Repositories.PullRequests.RequestChanges(opt)
	if err != nil {
		return nil, err
	}

	var participant Participant
	err = mapstructure.Decode(response, &participant)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (c Client) PrUnApprove(repoOrga string, repoSlug string, id string) error {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	_, err := client.Repositories.PullRequests.UnApprove(opt)
	return err
}

func (c Client) PrUnRequestChanges(repoOrga string, repoSlug string, id string) error {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		ID:       id,
	}

	_, err := client.Repositories.PullRequests.UnRequestChanges(opt)
	return err
}
