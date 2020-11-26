package repo

import (
	"github.com/craftamap/bb/cmd/commands/repo/view"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

func Add(rootCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	repoCommand := cobra.Command{
		Use:   "repo",
		Long:  "Work with repositories",
		Short: "Manage your repository",
	}

	view.Add(&repoCommand, globalOpts)

	rootCmd.AddCommand(&repoCommand)
}
