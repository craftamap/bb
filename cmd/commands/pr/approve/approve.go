package approve

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	Unapprove bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	approveCmd := &cobra.Command{
		Use:   "approve <number of id>",
		Short: "approve a pull request",
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
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
			} else {
				branchName, err := git.CurrentBranch()
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}

				prs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, branchName)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				if len(prs.Values) == 0 {
					fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Nothing on this branch")
					return
				}

				id = prs.Values[0].ID
			}
			if !Unapprove {
				participant, err := c.PrApprove(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				fmt.Println(aurora.Green(":: "), participant.State)
			} else {
				err := c.PrUnApprove(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
				if err != nil && !strings.Contains(err.Error(), "204") {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				fmt.Println(aurora.Green(":: "), "unapproved")
			}
		},
	}
	approveCmd.Flags().BoolVar(&Unapprove, "unapprove", false, "remove approval")
	prCmd.AddCommand(approveCmd)
}
