package internal

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ktrysmt/go-bitbucket"
)

func (c Client) DiffGet(repoOrga string, repoSlug string, spec string) (string, error) {
	client := bitbucket.NewBasicAuth(c.Username, c.Password)

	opt := bitbucket.DiffOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
		Spec:     spec,
	}

	responseBodyI, err := client.Repositories.Diff.GetDiff(&opt)
	if err != nil {
		return "", err
	}
	responseBody, ok := responseBodyI.(io.ReadCloser)
	if !ok {
		return "", fmt.Errorf("responseBody is no io.ReadCloser")
	}

	byteBody, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return "", err
	}

	return string(byteBody), nil
}
