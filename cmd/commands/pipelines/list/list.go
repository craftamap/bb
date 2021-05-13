package list

import (
	"fmt"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

const (
	SUCCESSFUL = "SUCCESSFUL"
	FAILED     = "FAILED"
	STOPPED    = "STOPPED"
	RUNNING    = "RUNNING"
)

func Add(pipelinesCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List pipeline executions this repository",
		Long:  "List pipeline executions this repository",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			pipelines, err := c.PipelineList(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				logging.Error(err)
				return
			}
			if len(*pipelines) == 0 {
				fmt.Println("No builds/statuses found for this pull request")
				return
			}
			hashToDescription := map[string]string{}

			for _, pipeline := range *pipelines {
				revHash := pipeline.Target.(map[string]interface{})["commit"].(map[string]interface{})["hash"].(string)
				commitDescription := ""
				commit, err := c.GetCommit(bbrepo.RepoOrga, bbrepo.RepoSlug, revHash)
				if err == nil {
					commitDescription = strings.Split(commit.Message, "\n")[0]
				}
				hashToDescription[revHash] = commitDescription
			}
			for _, pipeline := range *pipelines {
				var statusIcon string
				switch pipeline.PipelineState.Result.Name {
				case SUCCESSFUL:
					statusIcon = aurora.Green("✓").String()
				case FAILED, STOPPED:
					statusIcon = aurora.Red("X").String()
				}
				switch pipeline.PipelineState.Stage.Name {
				case RUNNING:
					statusIcon = aurora.Yellow("⏱️").String()
				}
				revHash := pipeline.Target.(map[string]interface{})["commit"].(map[string]interface{})["hash"].(string)
				description := hashToDescription[revHash]
				revHashShort := revHash[:7]

				branchName := pipeline.Target.(map[string]interface{})["ref_name"].(string)

				fmt.Printf(
					"%s #%d %s (%ds) %s\n",
					statusIcon,
					aurora.Index(242, pipeline.BuildNumber),
					description,
					pipeline.BuildSecondsUsed,
					aurora.Index(242,
						fmt.Sprintf("(%s, %s, %s)",
							pipeline.Creator.DisplayName,
							revHashShort,
							branchName,
						)),
				)
			}
		},
	}

	pipelinesCmd.AddCommand(listCmd)
}
