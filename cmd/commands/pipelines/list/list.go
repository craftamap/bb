package list

import (
	"github.com/craftamap/bb/cmd/options"
	"github.com/spf13/cobra"
)

func Add(pipelinesCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List and filter pull requests in this repository",
		Long:  "List and filter pull requests in this repository",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo
		},
	}

	pipelinesCmd.AddCommand(listCmd)
}
