package git

import (
	"fmt"
	"strings"

	"github.com/cli/cli/git"
)

type BitbucketRepo struct {
	RepoOrga string
	RepoSlug string
	Remote   git.Remote
}

func GetBitbucketRepo(remoteName string) (*BitbucketRepo, error) {
	remotes, err := git.Remotes()
	if err != nil {
		return nil, err
	}

	var selectedRemote git.Remote
	for _, remote := range remotes {
		if remote.Name == remoteName {
			selectedRemote = *remote
		}
	}
	// If no selectedRemote is found, throw an error
	if selectedRemote.Name == "" {
		return nil, fmt.Errorf("could not find the specified remote %s", remoteName)
	}

	path := strings.Split(selectedRemote.FetchURL.Path, "/")[1:]

	repoOrga := path[0]
	repoSlug := path[1]

	if selectedRemote.FetchURL.Scheme == "ssh" && strings.HasSuffix(repoSlug, ".git") {
		repoSlug = strings.TrimSuffix(repoSlug, ".git")
	}

	bbrepo := BitbucketRepo{
		Remote:   selectedRemote,
		RepoOrga: repoOrga,
		RepoSlug: repoSlug,
	}

	return &bbrepo, nil
}

func (b *BitbucketRepo) IsBitbucketOrg() bool {
	return strings.Contains(b.Remote.FetchURL.String(), "bitbucket")
}
