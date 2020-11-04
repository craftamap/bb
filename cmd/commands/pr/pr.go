package pr

import (
	"github.com/spf13/cobra"

	"github.com/craftamap/bb/cmd/commands/pr/create"
	"github.com/craftamap/bb/cmd/commands/pr/list"
	"github.com/craftamap/bb/cmd/commands/pr/merge"
	"github.com/craftamap/bb/cmd/commands/pr/statuses"
	"github.com/craftamap/bb/cmd/commands/pr/view"
	"github.com/craftamap/bb/cmd/options"
)

func Add(rootCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	prCommand := cobra.Command{
		Use:   "pr",
		Long:  "Work with pull requests",
		Short: "Manage pull requests",
	}

	list.Add(&prCommand, globalOpts)
	view.Add(&prCommand, globalOpts)
	create.Add(&prCommand, globalOpts)
	statuses.Add(&prCommand, globalOpts)
	merge.Add(&prCommand, globalOpts)

	rootCmd.AddCommand(&prCommand)
}
