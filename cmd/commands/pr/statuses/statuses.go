package statuses

import (
	"fmt"
	"github.com/craftamap/bb/util/logging"
	"strconv"
	"strings"

	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

const (
	SUCCESSFUL = "SUCCESSFUL"
	FAILED     = "FAILED"
	STOPPED    = "STOPPED"
	INPROGRESS = "INPROGRESS"
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	statusesCmd := &cobra.Command{
		Use:   "statuses [<number of pr>]",
		Short: "Show CI status for a single pull request",
		Long:  "Show CI status for a single pull request",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if len(args) > 0 {
				id, err = strconv.Atoi(strings.TrimPrefix(args[0], "#"))
				if err != nil {
					logging.Error(err)
					return
				}
			} else {
				branchName, err := git.CurrentBranch()
				if err != nil {
					logging.Error(err)
					return
				}

				prs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, branchName)
				if err != nil {
					logging.Error(err)
					return
				}
				if len(prs.Values) == 0 {
					logging.Warning("Nothing on this branch")
					return
				}

				id = prs.Values[0].ID
			}
			statuses, err := c.PrStatuses(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}

			if len(statuses.Values) == 0 {
				fmt.Println("No builds/statuses found for this pull request")
			} else {
				var (
					allChecksSuccessful = true
					successfulCount     = 0
					failedCount         = 0
					inProgressCount     = 0
					stoppedCount        = 0
				)

				for _, status := range statuses.Values {
					if status.State != SUCCESSFUL {
						allChecksSuccessful = false
					}

					switch status.State {
					case SUCCESSFUL:
						successfulCount++
					case FAILED:
						failedCount++
					case INPROGRESS:
						inProgressCount++
					case STOPPED:
						stoppedCount++
					}
				}
				if allChecksSuccessful {
					fmt.Println(aurora.Bold("All checks were successful").String())
				}
				fmt.Printf("%d failed, %d successful, %d in progress and %d stopped\n", failedCount, successfulCount, inProgressCount, stoppedCount)
				fmt.Println()

				for _, status := range statuses.Values {
					var statusIcon string
					switch status.State {
					case SUCCESSFUL:
						statusIcon = aurora.Green("✓").String()
					case FAILED, STOPPED:
						statusIcon = aurora.Red("X").String()
					case INPROGRESS:
						statusIcon = aurora.Yellow("⏱️").String()
					}

					fmt.Printf("%s %s %s %s\n", statusIcon, aurora.Index(242, status.Type), status.Name, aurora.Index(242, status.URL))
				}
			}
		},
	}

	prCmd.AddCommand(statusesCmd)
}
