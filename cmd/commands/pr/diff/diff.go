package diff

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	mergeCmd := &cobra.Command{
		Use:   "diff [<number of pr>]",
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
			repo, err := c.PrView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}

			diff, err := c.DiffGet(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%s..%s", repo.Source.Commit.Hash, repo.Destination.Commit.Hash))
			if err != nil {
				logging.Error(err)
				return
			}

			diffLines := bufio.NewScanner(strings.NewReader(diff))
			for diffLines.Scan() {
				diffLine := diffLines.Text()
				switch {
				case isHeaderLine(diffLine):
					fmt.Printf("\x1b[1;38m%s\x1b[m\n", diffLine)
				case isAdditionLine(diffLine):
					fmt.Printf("\x1b[32m%s\x1b[m\n", diffLine)
				case isRemovalLine(diffLine):
					fmt.Printf("\x1b[31m%s\x1b[m\n", diffLine)
				default:
					fmt.Println(diffLine)
				}
			}
		},
	}
	prCmd.AddCommand(mergeCmd)
}

var diffHeaderPrefixes = []string{"+++", "---", "diff", "index"}

func isHeaderLine(dl string) bool {
	for _, p := range diffHeaderPrefixes {
		if strings.HasPrefix(dl, p) {
			return true
		}
	}
	return false
}

func isAdditionLine(dl string) bool {
	return strings.HasPrefix(dl, "+")
}

func isRemovalLine(dl string) bool {
	return strings.HasPrefix(dl, "-")
}
