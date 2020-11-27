package list

import (
	"fmt"

	"github.com/craftamap/bb/cmd/options"
	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	Web bool
)

func Add(downloadsCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	listCmd := &cobra.Command{
		Use: "list",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if Web {
				err := browser.OpenURL(fmt.Sprintf("https://bitbucket.org/%s/%s/downloads", bbrepo.RepoOrga, bbrepo.RepoSlug))
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				return
			}
			downloads, err := c.GetDownloads(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
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

	listCmd.Flags().BoolVar(&Web, "web", false, "show the downloads in your browsers")

	downloadsCmd.AddCommand(listCmd)
}
