package delete

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/util/logging"
	"github.com/spf13/cobra"
)

func Add(issueCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	deleteCmd := &cobra.Command{
		Use:   "delete <nr of issue>",
		Short: "delete an issue",
		Long:  "delete an issue",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if len(args) > 0 {
				id, err = strconv.Atoi(strings.TrimPrefix(args[0], "#"))
				if err != nil {
					logging.Error(err)
					return
				}
			}

			response, err := c.IssuesDelete(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err, response)
				return
			}
			logging.Success(fmt.Sprintf("issue %d deleted", id))

		},
	}
	issueCmd.AddCommand(deleteCmd)
}
