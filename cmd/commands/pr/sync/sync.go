package sync

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cli/cli/git"
	"github.com/cli/safeexec"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/internal/run"
	"github.com/craftamap/bb/util/logging"
	"github.com/spf13/cobra"
)

var (
	Rebase bool
	Push   bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	syncCmd := &cobra.Command{
		Use:   "sync [<number of pr>]",
		Long:  "Sync a pull request on Bitbucket.org",
		Short: "Sync a pull request",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

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

			pr, err := c.PrView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}

			remoteDestinationBranch := fmt.Sprintf("origin/%s", pr.Destination.Branch.Name)

			var cmdQueue [][]string
			cmdQueue = append(cmdQueue, []string{"git", "fetch", "origin", pr.Destination.Branch.Name})
			if Rebase {
				cmdQueue = append(cmdQueue, []string{"git", "rebase", remoteDestinationBranch})
			} else {
				cmdQueue = append(cmdQueue, []string{"git", "merge", "--commit", remoteDestinationBranch})
			}
			if Push {
				cmdQueue = append(cmdQueue, []string{"git", "push"})
			}

			processCmdQueue(cmdQueue)

		},
	}

	syncCmd.Flags().BoolVarP(&Rebase, "rebase", "r", false, "uses rebase instead of merge to sync")
	syncCmd.Flags().BoolVar(&Push, "push", true, "push after merge/rebase")

	prCmd.AddCommand(syncCmd)
}

func processCmdQueue(cmdQueue [][]string) error {
	for _, args := range cmdQueue {
		exe, err := safeexec.LookPath(args[0])

		if err != nil {
			return err
		}
		cmd := exec.Command(exe, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := run.PrepareCmd(cmd).Run(); err != nil {
			return err
		}
	}
	return nil
}
