package pipelines

import (
	"github.com/craftamap/bb/cmd/commands/pipelines/list"
	"github.com/craftamap/bb/cmd/commands/pipelines/logs"
	"github.com/craftamap/bb/cmd/commands/pipelines/view"
	"github.com/craftamap/bb/cmd/commands/pipelines/wait"
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
	view.Add(&pipelineCommand, globalOpts)
	logs.Add(&pipelineCommand, globalOpts)
	wait.Add(&pipelineCommand, globalOpts)

	rootCmd.AddCommand(&pipelineCommand)
}
