package checkout

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/cli/cli/git"
	"github.com/cli/safeexec"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/internal/run"
	"github.com/spf13/cobra"
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	checkoutCmd := &cobra.Command{
		Use:   "checkout [<number of pr>]",
		Long:  "checkout a pull request on Bitbucket.org",
		Short: "checkout a pull request",
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
				//TODO: Error
			}
			pr, err := c.PrView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}

			newBranchName := pr.Source.Branch.Name
			// for now, we do not support prs from other repositories
			remoteBranch := fmt.Sprintf("origin/%s", newBranchName)

			var cmdQueue [][]string

			if _, err := git.ShowRefs("refs/heads/" + newBranchName); err == nil {
				cmdQueue = append(cmdQueue, []string{"git", "checkout", newBranchName})
				cmdQueue = append(cmdQueue, []string{"git", "merge", "--ff-only", fmt.Sprintf("refs/remotes/%s", remoteBranch)})
			} else {
				cmdQueue = append(cmdQueue, []string{"git", "checkout", "-b", newBranchName, "--no-track", remoteBranch})
				cmdQueue = append(cmdQueue, []string{"git", "config", fmt.Sprintf("branch.%s.remote", newBranchName), "origin"})
				cmdQueue = append(cmdQueue, []string{"git", "config", fmt.Sprintf("branch.%s.merge", newBranchName), "refs/heads/" + newBranchName})
			}

			for _, args := range cmdQueue {
				exe, err := safeexec.LookPath(args[0])

				if err != nil {
					logging.Error(err)
					return
				}
				cmd := exec.Command(exe, args[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := run.PrepareCmd(cmd).Run(); err != nil {
					logging.Error(err)
					return
				}
			}

		},
	}
	prCmd.AddCommand(checkoutCmd)
}
