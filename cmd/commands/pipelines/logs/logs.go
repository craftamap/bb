package logs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/craftamap/bb/util/logging"

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
			pipelineID, err := strconv.Atoi(strings.Replace(args[0], "#", "", 1))
			if err != nil {
				logging.Error(err)
				return
			}
			steps, err := c.PipelineStepsList(bbrepo.RepoOrga, bbrepo.RepoSlug, strconv.Itoa(pipelineID))
			if err != nil {
				logging.Error(err)
				return
			}

			for _, step := range *steps {
				logging.Success(fmt.Sprint(aurora.Bold("step name:"), step.Name))
				log, _ := c.PipelinesLogs(bbrepo.RepoOrga, bbrepo.RepoSlug, strconv.Itoa(pipelineID), step.UUID)
				fmt.Println(log)
			}
		},
	}
	pipelinesCmd.AddCommand(logsCmd)
}
