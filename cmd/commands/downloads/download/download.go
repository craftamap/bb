package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			downloadLink := dwnld.Links["self"].Href

			fmt.Printf("%s%s\n", aurora.Green(":: "), fmt.Sprintf("Downloading file from %s", downloadLink))

			req, err := http.NewRequest("GET", downloadLink, strings.NewReader(""))
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}
			req.SetBasicAuth(c.Username, c.Password)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("%s%s\n", aurora.Green(":: "), "Downloaded!")
			fmt.Printf("%s%s\n", aurora.Green(":: "), "Saving file to .")

			out, err := os.Create("./" + args[0])
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

		},
	}

	downloadCmd.Flags().BoolVar(&Web, "web", false, "view the pull request in your browser")
	prCmd.AddCommand(downloadCmd)
}

func downloadsToMap(downloads *internal.Downloads) map[string]internal.Download {
	downloadMap := map[string]internal.Download{}
	for _, dwnld := range downloads.Values {
		downloadMap[dwnld.Name] = dwnld
	}
	return downloadMap
}
