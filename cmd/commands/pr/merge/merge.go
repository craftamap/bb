package merge

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/craftamap/bb/client"
	"github.com/craftamap/bb/util/logging"

	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/commands/pr/view"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

var (
	CloseBranch bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	var mergeCmd *cobra.Command
	mergeCmd = &cobra.Command{
		Use:   "merge [<number of pr>]",
		Long:  "Merge a pull request on Bitbucket.org",
		Short: "Merge a pull request",
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
			} else if globalOpts.IsFSRepo {
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
			} else {
				logging.Warning("Not in a repository and no PR selected")
				return
			}

			var pr *client.PullRequest

			if mergeCmd.Flag("close-source-branch").Changed {
				pr, err = c.PrMergeWithCloseBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id), CloseBranch)
				if err != nil {
					logging.Error(err)
					return
				}
			} else {
				pr, err = c.PrMerge(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
				if err != nil {
					logging.Error(err)
					return
				}
			}

			commits, err := c.PrCommits(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}
			view.PrintSummary(pr, commits)
		},
	}

	mergeCmd.Flags().BoolVar(&CloseBranch, "close-source-branch", false, "close the source branch (pr setting if omitted)")

	prCmd.AddCommand(mergeCmd)
}
