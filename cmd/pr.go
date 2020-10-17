package cmd

import (
	"github.com/craftamap/bb/internal"
	"github.com/spf13/cobra"
)

var (
	prCommand = cobra.Command{
		Use: "pr",
		Run: func(cmd *cobra.Command, args []string) {
			internal.PrList(globalOpts.Username, globalOpts.Password, globalOpts.RepoOrga, globalOpts.RepoSlug)
		},
	}
	prListCommand = cobra.Command{
		Use: "list",
		Run: list,
	}
	prViewCommand = cobra.Command{
		Use: "view",
		Run: view,
	}
	prCreateCommand = cobra.Command{
		Use: "create",
		Run: create,
	}
)

func init() {
	prCommand.AddCommand(&prListCommand)
	prCommand.AddCommand(&prViewCommand)
	prCommand.AddCommand(&prCreateCommand)
	rootCmd.AddCommand(&prCommand)
}

func list(cmd *cobra.Command, args []string) {
}

func view(cmd *cobra.Command, args []string) {
}

func create(cmd *cobra.Command, args []string) {

}
