package internal

import (
	"strings"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

type Project struct {
	Key  string
	Name string
}

type Repository struct {
	Links       map[string]interface{} `mapstructure:"Links"`
	UUID        string                 `mapstructure:"Uuid"`
	FullName    string                 `mapstructure:"Full_name"`
	IsPrivate   bool                   `mapstructure:"Is_private"`
	Owner       *Account               `mapstructure:"Owner"`
	Name        string                 `mapstructure:"Name"`
	Description string                 `mapstructure:"Description"`
	Size        int                    `mapstructure:"Size"`
	Language    string                 `mapstructure:"Language"`
	HasIssues   bool                   `mapstructure:"Has_issues"`
	ForkPolicy  string                 `mapstructure:"ForkPolicy"`
	MainBranch  *Branch                `mapstructure:"Mainbranch"`
	Project     Project                `mapstructure:"Project"`
	// Parent      *Repository            `mapstructure:"Parent"`
	// CreatedOn   time.Time              `mapstructure:"CreatedOn"`
	// UpdatedOn   time.Time              `mapstructure:"UpdatedOn"`
}

type DefaultReviewers struct {
	Values []*Account `json:"values"`
}

func (c Client) RepositoryGet(repoOrga string, repoSlug string) (*Repository, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := &bitbucket.RepositoryOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
	}

	repositoryResponse, err := client.Repositories.Repository.Get(opt)
	if err != nil {
		return nil, err
	}

	var repo Repository
	err = mapstructure.Decode(repositoryResponse, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (c Client) GetDefaultReviewers(repoOrga string, repoSlug string) (*DefaultReviewers, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)
	opt := &bitbucket.RepositoryOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
	}

	response, err := client.Repositories.Repository.ListDefaultReviewers(opt)
	if err != nil {
		return nil, err
	}

	defaultReviewers := DefaultReviewers{}
	err = mapstructure.Decode(response, &defaultReviewers)
	if err != nil {
		return nil, err
	}

	return &defaultReviewers, nil
}

func (c Client) GetReadmeContent(repoOrga string, repoSlug string, branch string) (string, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	allRootFiles, err := client.Repositories.Repository.ListFiles(&bitbucket.RepositoryFilesOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		Ref:      branch,
	})
	if err != nil {
		return "", err
	}

	potentialReadmes := []bitbucket.RepositoryFile{}

	for _, rootFile := range allRootFiles {
		if strings.HasPrefix(strings.ToLower(rootFile.Path), "readme") {
			potentialReadmes = append(potentialReadmes, rootFile)
		}
	}

	if len(potentialReadmes) <= 0 {
		return "", nil
	}

	fileName := potentialReadmes[0].Path

	fileBlob, err := client.Repositories.Repository.GetFileBlob(&bitbucket.RepositoryBlobOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		Ref:      branch,
		Path:     fileName,
	})
	if err != nil {
		return "", err
	}

	return fileBlob.String(), nil
}
