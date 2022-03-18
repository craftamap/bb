package git

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cli/cli/git"
	"github.com/craftamap/bb/internal/run"
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

func CurrentHead() (string, error) {
	headCmd, err := git.GitCommand("rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	output, err := run.PrepareCmd(headCmd).Output()
	return firstLine(output), err
}

func RepoPath() (string, error) {
	pathCmd, err := git.GitCommand("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	output, err := run.PrepareCmd(pathCmd).Output()
	return firstLine(output), err
}

func firstLine(output []byte) string {
	if i := bytes.IndexAny(output, "\n"); i >= 0 {
		return string(output)[0:i]
	}
	return string(output)
}
