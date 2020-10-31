package downloads

import (
	"github.com/craftamap/bb/cmd/commands/downloads/list"
	"github.com/craftamap/bb/cmd/commands/downloads/upload"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

func Add(rootCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	downloadsCmd := &cobra.Command{
		Use:   "downloads",
		Short: "Manage repository downloads",
		Long:  "Manage repository downloads on Bitbucket.org",
	}

	list.Add(downloadsCmd, globalOpts)
	upload.Add(downloadsCmd, globalOpts)

	rootCmd.AddCommand(downloadsCmd)
}
