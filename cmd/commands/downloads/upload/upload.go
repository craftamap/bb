package upload

import (
	"fmt"
	"github.com/craftamap/bb/util/logging"
	"os"
	"path/filepath"

	"github.com/craftamap/bb/cmd/options"
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
				logging.Warning("No file specified")
				return
			}

			fpath := args[0]
			fmt.Println(fpath)

			if _, err := os.Stat(fpath); os.IsNotExist(err) {
				logging.Error(err)
				return
			}
			logging.Success(fmt.Sprintf("Uploading file %s", filepath.Base(fpath)))

			_, err := c.UploadDownload(bbrepo.RepoOrga, bbrepo.RepoSlug, fpath)
			if err != nil {
				logging.Error(err)
				return
			}

			//if err != nil {
			//	logging.Error(err)
			//	return
			//}

			logging.Success(fmt.Sprintf("Uploaded file %s\n", filepath.Base(fpath)))
		},
	}

	downloadsCmd.AddCommand(uploadCmd)
}
