package download

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
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
			var (
				remoteName  string
				storagePath string
			)

			if len(args) == 0 {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "No file selected")
				return
			}
			if len(args) >= 1 {
				remoteName = args[0]
				storagePath = "."
			}
			if len(args) == 2 {
				storagePath = args[1]
			}
			if len(args) > 2 {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "Too many arguments")
				return
			}

			// First, check if dir exists and is dir. if not, fail.
			// Then check if filename exists. If not, use it as filename
			// Else, check if its file or dir
			//   If file, use it as filename
			//   Else, use storagePath as dir and remoteName as file

			dir, fname := filepath.Split(storagePath)
			if dir != "" {
				info, err := os.Stat(dir)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				if !info.IsDir() {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), fmt.Sprintf("%s is not a directory", dir))
				}
			}

			info, err := os.Stat(storagePath)
			if !os.IsNotExist(err) && info.IsDir() {
				dir = storagePath
				_, fname = filepath.Split(remoteName)
			}
			storagePath = filepath.Join(dir, fname)

			// Safe downloads collition-free
			info, err = os.Stat(storagePath)
			if !os.IsNotExist(err) && !info.IsDir() {
				dir, file := filepath.Split(storagePath)
				file += "." + RandomString(10)
				storagePath = filepath.Join(dir, file)
			}

			// Actual logic
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			fmt.Printf("%s%s\n", aurora.Green(":: "), "Getting all downloads")

			downloads, err := c.GetDownloads(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			downloadMap := downloadsToMap(downloads)
			dwnld, ok := downloadMap[remoteName]
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
			fmt.Printf("%s%s\n", aurora.Green(":: "), fmt.Sprintf("Saving file to %s", storagePath))

			out, err := os.Create(storagePath)
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

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		a, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		s[i] = letters[a.Int64()]
	}
	return string(s)
}
