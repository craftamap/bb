package list

import (
	"fmt"
	"strings"

	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/util/logging"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	Web        bool
	States     []string
	Types      []string
	Priorities []string
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

				linkBase := repo.Links["html"].(map[string]interface{})["href"].(string)
				link := linkBase + "/issues"
				err = browser.OpenURL(link)
				if err != nil {
					logging.Error(err)
					return
				}

				return
			}

			for _, state := range States {
				if strings.ToUpper(state) == "ALL" {
					States = []string{"new", "open", "resolved", "on hold", "invalid", "duplicate", "wontfix", "closed"}
				}
			}

			for _, typus := range Types {
				if strings.ToUpper(typus) == "ALL" {
					Types = []string{"bug", "enhancement", "proposal", "task"}
				}
			}

			for _, priority := range Priorities {
				if strings.ToUpper(priority) == "ALL" {
					Priorities = []string{"trivial", "minor", "major", "critical", "blocker"}
				}
			}

			issues, err := c.IssuesList(bbrepo.RepoOrga, bbrepo.RepoSlug, States, Types, Priorities)
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
				fmt.Printf("#%03d %s  %s   %s\n", aurora.Green(issue.ID), state, issue.Title, aurora.Index(242, fmt.Sprintf("by %s", issue.Reporter.DisplayName)))
			}
		},
	}
	listCmd.Flags().StringArrayVar(&States, "state", []string{"new", "open"}, "Filter by state: {new|open|resolved|on hold|invalid|duplicate|wontfix|closed|all}")
	listCmd.Flags().StringArrayVar(&Types, "type", []string{"all"}, "Filter by type/kind: {bug|enhancement|proposal|task|all}")
	listCmd.Flags().StringArrayVar(&Priorities, "priority", []string{"all"}, "Filter by priority: {trivial|minor|major|critical|blocker}")
	listCmd.Flags().BoolVar(&Web, "web", false, "view issues in your browser")
	issueCmd.AddCommand(listCmd)
}
