package issue

import (
	"github.com/craftamap/bb/cmd/commands/issue/create"
	"github.com/craftamap/bb/cmd/commands/issue/list"
	"github.com/craftamap/bb/cmd/commands/issue/view"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

func Add(rootCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	issueCommand := cobra.Command{
		Use:   "issue",
		Long:  "Work with issues",
		Short: "Manage issues",
	}

	list.Add(&issueCommand, globalOpts)
	view.Add(&issueCommand, globalOpts)
	create.Add(&issueCommand, globalOpts)

	rootCmd.AddCommand(&issueCommand)
}
