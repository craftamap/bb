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
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	Rebase bool
	Push   bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Long:  "synchronizes the current pull request with new changes in it's destination branch by either merging or rebasing the changes locally. If the rebase or merge fails because of an conflict, the merge must be resolved manually.",
		Short: "Sync a pull request locally",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(_ *cobra.Command, _ []string) {
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

			err = processCmdQueue(cmdQueue)
			if err != nil {
				logging.Error(err)
				logging.Note("Looks like an error occured - you probably need to resolve the conflict manually now. Start by running " + aurora.BgBrightBlack(aurora.White(" git status ")).String() + " to get hints on how to resolve the conflicts.")
				cancelCommand := aurora.BgBrightBlack(aurora.White(" git merge --abort ")).String()
				if Rebase {
					cancelCommand = aurora.BgBrightBlack(aurora.White(" git rebase --abort ")).String()
				}
				logging.Note("If you don't know how to resolve a conflict manually, you can run " + cancelCommand + " to reset.")
				return
			}

			if Push {
				cmdQueue = [][]string{{"git", "push"}}
				err = processCmdQueue(cmdQueue)
				if err != nil {
					logging.Error(err)
				}
			}
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
