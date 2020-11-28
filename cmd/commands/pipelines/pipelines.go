package pipelines

import (
	"github.com/craftamap/bb/cmd/commands/pipelines/list"
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

func Add(rootCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	pipelineCommand := cobra.Command{
		Use:   "pipelines",
		Long:  "Work with pipelines",
		Short: "Manage pull pipelines",
	}

	list.Add(&pipelineCommand, globalOpts)

	rootCmd.AddCommand(&pipelineCommand)
}
