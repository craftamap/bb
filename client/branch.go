package client

import "github.com/ktrysmt/go-bitbucket"

func (c Client) GetBranch(repoOrga string, repoSlug string, branchName string) (*bitbucket.RepositoryBranch, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)
	return client.Repositories.Repository.GetBranch(&bitbucket.RepositoryBranchOptions{
		Owner:      repoOrga,
		RepoSlug:   repoSlug,
		BranchName: branchName,
	})
}
