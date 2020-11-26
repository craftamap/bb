package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func Add(downloadsCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	uploadCmd := &cobra.Command{
		Use: "upload",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if len(args) == 0 {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "No file specified")
				return
			}

			fpath := args[0]
			fmt.Println(fpath)

			if _, err := os.Stat(fpath); os.IsNotExist(err) {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}
			fmt.Printf("%s Uploading file %s\n", aurora.Green(":: "), filepath.Base(fpath))

			_, err := c.UploadDownload(bbrepo.RepoOrga, bbrepo.RepoSlug, fpath)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			//if err != nil {
			//	fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
			//	return
			//}

			fmt.Printf("%s Uploaded file %s\n", aurora.Green(":: "), filepath.Base(fpath))
		},
	}

	downloadsCmd.AddCommand(uploadCmd)
}
