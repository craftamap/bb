package internal

import (
	"errors"
	"fmt"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/logrusorgru/aurora"
)

func PrList(username string, password string, repoOrga string, repoSlug string) error {

	client := bitbucket.NewBasicAuth(username, password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
	}

	response, err := client.Repositories.PullRequests.Gets(opt)
	if err != nil {
		return err
	}
	mapResponse, ok := response.(map[string]interface{})
	if !ok {
		return errors.New("type assertion failed")
	}
	values, ok := mapResponse["values"].([]interface{})
	if !ok {
		return errors.New("type assertion failed")
	}

	fmt.Println()
	fmt.Printf("%s Showing %d of %d open pull requests in %s/%s\n", aurora.Blue(" :: "), len(values), int(mapResponse["size"].(float64)), repoOrga, repoSlug)
	fmt.Println()
	for _, pr := range values {
		mapPr, ok := pr.(map[string]interface{})
		if !ok {
			return errors.New("type assertion failed")
		}
		source, ok := mapPr["source"].(map[string]interface{})
		if !ok {
			return errors.New("type assertion failed")
		}

		sourceBranch, ok := source["branch"].(map[string]interface{})
		if !ok {
			return errors.New("type assertion failed")
		}

		sourceName := sourceBranch["name"]

		dest, ok := mapPr["destination"].(map[string]interface{})
		if !ok {
			return errors.New("type assertion failed")
		}

		destBranch, ok := dest["branch"].(map[string]interface{})
		if !ok {
			return errors.New("type assertion failed")
		}

		destName := destBranch["name"]
		fmt.Printf("#%03d  %s   %s -> %s\n", aurora.Green(int(mapPr["id"].(float64))), mapPr["title"].(string), sourceName, destName)
	}

	return nil
}
