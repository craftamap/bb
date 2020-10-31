package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func Add(downloadsCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	uploadCmd := &cobra.Command{
		Use: "upload",
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

			if len(args) == 0 {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "No file specified")
				return
			}

			fpath := args[0]
			fmt.Println(fpath)

			if _, err := os.Stat(fpath); os.IsNotExist(err) {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			fmt.Printf("%s Uploading file %s\n", aurora.Green(":: "), filepath.Base(fpath))

			// Workaround: As UploadDownload files currently, we need to recover
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println(aurora.Bold("WORKAROUND"), ": Recovered", r)
					}
				}()
				c.UploadDownload(bbrepo.RepoOrga, bbrepo.RepoSlug, fpath)
			}()

			//if err != nil {
			//	fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
			//	return
			//}

			fmt.Printf("%s Uploaded file %s\n", aurora.Green(":: "), filepath.Base(fpath))
		},
	}

	downloadsCmd.AddCommand(uploadCmd)
}
