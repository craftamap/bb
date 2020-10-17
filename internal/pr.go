package internal

import (
	"fmt"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/logrusorgru/aurora"
)

func PrList(username string, password string, repoOrga string, repoSlug string) {
	client := bitbucket.NewBasicAuth(username, password)

	opt := &bitbucket.PullRequestsOptions{
		Owner:    repoOrga,
		RepoSlug: repoSlug,
	}

	response, _ := client.Repositories.PullRequests.Gets(opt)
	mapResponse := response.(map[string]interface{})
	values := mapResponse["values"].([]interface{})

	fmt.Println()
	fmt.Printf("%s Showing %d of %d open pull requests in %s/%s\n", aurora.Blue(" :: "), len(values), int(mapResponse["size"].(float64)), repoOrga, repoSlug)
	fmt.Println()
	for _, pr := range values {
		mapPr := pr.(map[string]interface{})
		sourceName := mapPr["source"].(map[string]interface{})["branch"].(map[string]interface{})["name"]
		destName := mapPr["destination"].(map[string]interface{})["branch"].(map[string]interface{})["name"]
		fmt.Printf("#%03d  %s   %s -> %s\n", aurora.Green(int(mapPr["id"].(float64))), mapPr["title"].(string), sourceName, destName)
	}
}
