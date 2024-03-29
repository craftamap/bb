package requestchanges

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

var (
	UnRequest bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	requestCmd := &cobra.Command{
		Use:   "request-changes <number of id>",
		Short: "request-changes on a pull request",
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

			if !UnRequest {
				participant, err := c.PrRequestChanges(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
				if err != nil {
					logging.Error(err)
					return
				}
				logging.Success(participant.State)
			} else {
				err := c.PrUnRequestChanges(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
				if err != nil && !strings.Contains(err.Error(), "204") {
					logging.Error(err)
					return
				}
				logging.Success("unrequested")
			}
		},
	}
	requestCmd.Flags().BoolVar(&UnRequest, "unrequest", false, "remove request for changes")
	prCmd.AddCommand(requestCmd)
}
