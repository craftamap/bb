package view

import (
	"fmt"
	"github.com/craftamap/bb/util/logging"
	"strconv"
	"strings"

	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

const (
	SUCCESSFUL  = "SUCCESSFUL"
	FAILED      = "FAILED"
	STOPPED     = "STOPPED"
	RUNNING     = "RUNNING"
	IN_PROGRESS = "IN_PROGRESS"
)

func Add(pipelinesCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	viewCmd := &cobra.Command{
		Use:   "view <number of pipeline>",
		Short: "view a specific pipeline execution",
		Long:  "view a specific pipeline execution",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if len(args) == 0 {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "Missing argument <number of pipeline>")
				return
			}
			pipelineID, err := strconv.Atoi(strings.Replace(args[0], "#", "", 1))
			if err != nil {
				logging.Error(err)
				return
			}

			pipeline, err := c.PipelineGet(bbrepo.RepoOrga, bbrepo.RepoSlug, strconv.Itoa(pipelineID))
			if err != nil {
				logging.Error(err)
				return
			}

			steps, err := c.PipelineStepsList(bbrepo.RepoOrga, bbrepo.RepoSlug, strconv.Itoa(pipelineID))
			if err != nil {
				logging.Error(err)
				return
			}

			var statusIcon fmt.Stringer
			switch pipeline.PipelineState.Result.Name {
			case SUCCESSFUL:
				statusIcon = aurora.Green("✓" + " " + pipeline.PipelineState.Result.Name)
			case FAILED, STOPPED:
				statusIcon = aurora.Red("X" + " " + pipeline.PipelineState.Result.Name)
			}
			switch pipeline.PipelineState.Stage.Name {
			case RUNNING:
				statusIcon = aurora.Yellow("⏱️" + pipeline.PipelineState.Stage.Name)
			}

			revHash := pipeline.Target.(map[string]interface{})["commit"].(map[string]interface{})["hash"].(string)
			revHashShort := revHash[:7]

			commitDescription := ""
			commit, err := c.GetCommit(bbrepo.RepoOrga, bbrepo.RepoSlug, revHash)
			if err == nil {
				commitDescription = strings.Split(commit.Message, "\n")[0]
			}
			description := commitDescription

			branchName := pipeline.Target.(map[string]interface{})["ref_name"].(string)

			fmt.Printf(
				"#%d %s\n",
				aurora.Index(242, pipeline.BuildNumber),
				description,
			)
			fmt.Println(statusIcon)
			fmt.Printf("Creator: %s\n", pipeline.Creator.DisplayName)
			fmt.Printf("Target: %s, %s\n",
				branchName,
				revHashShort,
			)

			if len(*steps) > 0 {
				fmt.Println("Steps:")
			}
			for _, step := range *steps {
				var statusIcon fmt.Stringer
				switch step.State.Result.Name {
				case SUCCESSFUL:
					statusIcon = aurora.Green("✓")
				case FAILED, STOPPED:
					statusIcon = aurora.Red("X")
				}
				switch step.State.Name {
				case IN_PROGRESS:
					statusIcon = aurora.Yellow("⏱️")
				}
				fmt.Printf("- %s %s (%s)\n", statusIcon, step.Name, step.UUID)
			}
		},
	}
	pipelinesCmd.AddCommand(viewCmd)
}
