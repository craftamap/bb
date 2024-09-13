package wait

import (
	"strconv"
	"strings"
	"time"

	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/util/logging"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

const (
	SUCCESSFUL  = "SUCCESSFUL"
	FAILED      = "FAILED"
	STOPPED     = "STOPPED"
	RUNNING     = "RUNNING"
	IN_PROGRESS = "IN_PROGRESS"
)

func Add(pipelinesCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	waitCmd := &cobra.Command{
		Use:   "wait <number of pipeline>",
		Short: "wait until a specific pipeline execution has finished",
		Long:  "wait until a specific pipeline execution has finished. On it's own, this command, this is not particulary useful, but it can be combined with other commands, like notify-send, to send you a notification when a pipeline is done.",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if len(args) == 0 {
				logging.Error("Missing argument <number of pipeline>")
				return
			}
			pipelineID, err := strconv.Atoi(strings.Replace(args[0], "#", "", 1))
			if err != nil {
				logging.Error(err)
				return
			}

			for {
				pipeline, err := c.PipelineGet(bbrepo.RepoOrga, bbrepo.RepoSlug, strconv.Itoa(pipelineID))
				if err != nil {
					logging.Error(err)
					return
				}
				if verbose {
					logging.Note("Current pipeline state: ", pipeline.PipelineState.Name)
				}

				if pipeline.PipelineState.Name == "COMPLETED" {
					break
				}
				time.Sleep(1 * time.Second)
			}

		},
	}
	waitCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbosity of the command; logs the current state")
	pipelinesCmd.AddCommand(waitCmd)
}
