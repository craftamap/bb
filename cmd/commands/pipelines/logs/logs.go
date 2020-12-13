package logs

import (
	"fmt"

	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func Add(pipelinesCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	logsCmd := &cobra.Command{
		Use:   "logs <number of pipeline>",
		Short: "view the logs of a specific pipeline execution",
		Long:  "view the logs of a specific pipeline execution",
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
			steps, err := c.PipelineStepsList(bbrepo.RepoOrga, bbrepo.RepoSlug, args[0])
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			for _, step := range *steps {
				fmt.Println(aurora.Green("::"), aurora.Bold("step name:"), step.Name)
				log, _ := c.PipelinesLogs(bbrepo.RepoOrga, bbrepo.RepoSlug, args[0], step.UUID)
				fmt.Println(log)
			}
		},
	}
	pipelinesCmd.AddCommand(logsCmd)
}
