package list

import (
	"fmt"

	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/util/logging"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/wbrefvem/go-bitbucket"
)

var (
	Web    bool
	States []string
)

func Add(issueCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List and filter issues in this repository",
		Long:  "List and filter issues in this repository",
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

				linkWrapper := repo.Links["Html"].(*bitbucket.SubjectTypesRepositoryEvents)
				link := linkWrapper.Href + "/issues"
				err = browser.OpenURL(link)
				if err != nil {
					logging.Error(err)
					return
				}

				return
			}

			issues, err := c.IssuesList(bbrepo.RepoOrga, bbrepo.RepoSlug, States)
			if err != nil {
				logging.Error(err)
			}

			fmt.Println()
			fmt.Printf("%sShowing %d of %d issues in %s/%s\n", aurora.Blue(" :: "), len(issues.Values), issues.Size, bbrepo.RepoOrga, bbrepo.RepoSlug)
			fmt.Println()
			for _, issue := range issues.Values {
				var state string
				switch issue.State {
				case "new":
					state = aurora.BgIndex(55, " NEW ").String() + "     "
				case "open":
					state = aurora.BgGray(12, " OPEN ").String() + "    "
				case "on hold":
					state = aurora.BgBlue(" ON HOLD ").String() + "  "
				case "invalid":
					state = aurora.BgRed(" INVALID ").String() + "  "
				case "resolved":
					state = aurora.BgGreen(" RESOLVED ").String()
				case "duplicate":
					state = aurora.BgYellow(" DUPLICATE ").String()
				case "wontfix":
					state = aurora.BgRed(" WONTFIX ").String() + "  "
				case "closed":
					state = aurora.BgGreen(" CLOSED ").String() + "   "
				default:
					state = issue.State
				}
				fmt.Printf("#%03d %s  %s\n", aurora.Green(issue.ID), state, issue.Title)
			}
		},
	}
	listCmd.Flags().StringArrayVar(&States, "states", []string{"new", "open", "resolved", "invalid", "duplicate", "wontfix"}, "Filter by state: {new|open|resolved|on hold|invalid|duplicate|wontfix|closed}")
	listCmd.Flags().BoolVar(&Web, "web", false, "view issues in your browser")
	issueCmd.AddCommand(listCmd)
}
