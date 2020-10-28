package statuses

import (
	"fmt"
	"strconv"

	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	statusesCmd := &cobra.Command{
		Use: "statuses",
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error

			c := internal.Client{
				Username: globalOpts.Username,
				Password: globalOpts.Password,
			}

			bbrepo, err := bbgit.GetBitbucketRepo()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			if !bbrepo.IsBitbucketOrg() {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Are you sure this is a bitbucket repo?")
				return
			}

			if len(args) > 0 {
				id, err = strconv.Atoi(args[0])
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}
			} else {
				branchName, err := git.CurrentBranch()
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}

				prs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, branchName)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}
				if len(prs.Values) == 0 {
					fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Nothing on this branch")
					return
				}

				id = prs.Values[0].ID
			}
			statuses, err := c.PrStatuses(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			if len(statuses.Values) == 0 {
				fmt.Printf("No builds/statuses found for this pull request")
			} else {
				var (
					allChecksSuccessful = true
					successfulCount     = 0
					failedCount         = 0
					inProgressCount     = 0
					stoppedCount        = 0
				)

				for _, status := range statuses.Values {
					if status.State != "SUCCESSFUL" {
						allChecksSuccessful = false
					}

					switch status.State {
					case "SUCCESSFUL":
						successfulCount++
					case "FAILED":
						failedCount++
					case "INPROGRESS":
						inProgressCount++
					case "STOPPED":
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
					case "SUCCESSFUL":
						statusIcon = aurora.Green("✓").String()
					case "FAILED", "STOPPED":
						statusIcon = aurora.Red("X").String()
					case "INPROGRESS":
						statusIcon = aurora.Yellow("⏱️").String()
					}

					fmt.Printf("%s %s %s %s\n", statusIcon, aurora.BrightBlack(status.Type), status.Name, aurora.BrightBlack(status.URL))
				}

			}

		},
	}

	prCmd.AddCommand(statusesCmd)
}
