package list

import (
	"fmt"

	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func Add(downloadsCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	listCmd := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
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
			downloads, err := c.GetDownloads(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			for i := len(downloads.Values) - 1; i >= 0; i-- {
				download := downloads.Values[i]
				fmt.Printf(
					"• %s - %s - %s - %s - %s\n",
					download.Name,
					aurora.Index(242, humanize.IBytes(uint64(download.Size))),
					aurora.Index(242, download.User.DisplayName),
					aurora.Index(242, fmt.Sprintf("%d Downloads", download.Downloads)),
					aurora.Index(242, download.CreatedOn),
				)
			}
		},
	}

	downloadsCmd.AddCommand(listCmd)
}
