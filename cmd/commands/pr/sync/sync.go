package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cli/cli/git"
	"github.com/cli/cli/utils"
	"github.com/cli/safeexec"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/internal/run"
	"github.com/craftamap/bb/util/logging"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Method string
	Push   bool
)

const (
	MethodOptionMerge  = "merge"
	MethodOptionRebase = "rebase"
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Long:  "synchronizes the current pull request with new changes in it's destination branch by either merging or rebasing the changes locally. If rebasing or merging fails because of an conflict, the merge must be resolved manually.",
		Short: "Sync a pull request locally",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			// In order to check if the method exists in the config, we need to check here
			syncMethodIsSet := viper.IsSet("sync-method")
			if !syncMethodIsSet {
				logging.Note(
					"You can configure your preferred way of syncing by adding the following line to your configuration: ",
					aurora.BgBrightBlack(aurora.White(" sync-method = merge ")).String(),
					" or ",
					aurora.BgBrightBlack(aurora.White(" sync-method = rebase ")).String(),
				)
			}
			if syncMethodIsSet && !cmd.Flags().Lookup("method").Changed {
				Method = viper.GetString("sync-method")
			}

			if Method != MethodOptionRebase && Method != MethodOptionMerge {
				return fmt.Errorf("\"%s\" is not a valid method (select one of these: %s, %s)", Method, MethodOptionRebase, MethodOptionMerge)
			}
			return nil
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
				logging.Error("No pull request on this branch")
				return
			}
			id = prs.Values[0].ID

			pr, err := c.PrView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}
			if ucc, err := git.UncommittedChangeCount(); err == nil && ucc > 0 {
				logging.Warning(utils.Pluralize(ucc, "uncommitted change"), "; This might lead to issues")
			}

			remoteName := bbrepo.Remote.Name

			remoteDestinationBranch := fmt.Sprintf("%s/%s", remoteName, pr.Destination.Branch.Name)

			var cmdQueue [][]string
			cmdQueue = append(cmdQueue, []string{"git", "fetch", remoteName, pr.Destination.Branch.Name})
			if Method == MethodOptionRebase {
				cmdQueue = append(cmdQueue, []string{"git", "rebase", remoteDestinationBranch})
			} else {
				cmdQueue = append(cmdQueue, []string{"git", "merge", "--commit", remoteDestinationBranch})
			}

			var builder strings.Builder
			for _, cmd := range cmdQueue {
				builder.WriteString("	")
				for _, cmdPart := range cmd {
					builder.WriteString(cmdPart)
					builder.WriteString(" ")
				}
				builder.WriteString("\n")
			}
			logging.Note(fmt.Sprintf("Syncing by running the following commands: \n%v", builder.String()))

			err = processCmdQueue(cmdQueue)
			if err != nil {
				logging.Error(err)
				logging.Note("Looks like an error occurred - you probably need to resolve the conflict manually now. Start by running " + aurora.BgBrightBlack(aurora.White(" git status ")).String() + " to get hints on how to resolve the conflicts.")
				cancelCommand := aurora.BgBrightBlack(aurora.White(" git merge --abort ")).String()
				if Method == MethodOptionRebase {
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

	syncCmd.Flags().StringVar(&Method, "method", "merge", "sync using merge or rebase (merge/rebase)")
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
