package download

import (
	"fmt"

	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	Web bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	downloadCmd := &cobra.Command{
		Use:   "download <filename> [<target path>]",
		Short: "download files from bitbucket",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Get rid of this "list" workaround

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			downloads, err := c.GetDownloads(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			downloadMap := downloadsToMap(downloads)

			dwnld, ok := downloadMap[args[0]]
			if !ok {
				fmt.Println("Zu hülfe!!")
				return
			}

			fmt.Println(dwnld.Links)
		},
	}

	downloadCmd.Flags().BoolVar(&Web, "web", false, "view the pull request in your browser")
	prCmd.AddCommand(downloadCmd)
}

func downloadsToMap(downloads *internal.Downloads) map[string]*internal.Download {
	downloadMap := map[string]*internal.Download{}
	for _, dwnld := range downloads.Values {
		downloadMap[dwnld.Name] = &dwnld
	}
	return downloadMap
}
