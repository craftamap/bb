package list

import (
	"fmt"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	Web   bool
	State string
)

type LinkWrapper struct {
	Href string `mapstructure:"href"`
}

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List and filter pull requests in this repository",
		Long:    "List and filter pull requests in this repository",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if Web {
				repo, err := c.RepositoryGet(bbrepo.RepoOrga, bbrepo.RepoSlug)
				if err != nil {
					logging.Error(err)
				}

				linkWrapper := repo.Links["Html"].(map[string]any)
				link := linkWrapper["Href"].(string) + "/pull-requests"
				err = browser.OpenURL(link)
				if err != nil {
					logging.Error(err)
					return
				}

				return
			}

			state := strings.ToUpper(State)

			prs, err := c.PrList(bbrepo.RepoOrga, bbrepo.RepoSlug, []string{state})
			if err != nil {
				logging.Error(err)
			}

			fmt.Println()
			fmt.Printf("%sShowing %d of %d open pull requests in %s/%s\n", aurora.Blue(" :: "), len(prs.Values), prs.Size, bbrepo.RepoOrga, bbrepo.RepoSlug)
			fmt.Println()
			for _, pr := range prs.Values {
				fmt.Printf("#%03d  %s   %s -> %s\n", aurora.Green(pr.ID), pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
			}
		},
	}
	listCmd.Flags().StringVar(&State, "state", "open", "Filter by state: {open|merged|declined|superseded}")
	listCmd.Flags().BoolVar(&Web, "web", false, "view pull requests in your browser")
	prCmd.AddCommand(listCmd)
}
