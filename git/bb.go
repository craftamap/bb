package git

import (
	"strings"

	"github.com/cli/cli/git"
)

type BitbucketRepo struct {
	RepoOrga string
	RepoSlug string
	Remote   git.Remote
}

func GetBitbucketRepo() (*BitbucketRepo, error) {
	remotes, err := git.Remotes()
	if err != nil {
		return nil, err
	}

	var origin git.Remote
	for _, remote := range remotes {
		if remote.Name == "origin" {
			origin = *remote
		}
	}

	path := strings.Split(origin.FetchURL.Path, "/")[1:]

	repoOrga := path[0]
	repoSlug := path[1]

	if origin.FetchURL.Scheme == "ssh" && strings.HasSuffix(repoSlug, ".git") {
		repoSlug = strings.TrimSuffix(repoSlug, ".git")
	}

	bbrepo := BitbucketRepo{
		Remote:   origin,
		RepoOrga: repoOrga,
		RepoSlug: repoSlug,
	}

	return &bbrepo, nil
}
