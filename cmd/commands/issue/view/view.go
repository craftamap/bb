package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/craftamap/bb/client"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/util/logging"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	Web bool
)

func Add(issueCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	viewCmd := &cobra.Command{
		Use:   "view [<nr of issue>]",
		Short: "View a issue",
		Long:  "Display the title, body, and other information about a issue.",
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

			issue, err := c.IssuesView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}
			if Web {
				err := browser.OpenURL(issue.Links["html"].Href)
				if err != nil {
					logging.Error(err)
					return
				}
				return
			}

			PrintSummary(issue)
		},
	}
	viewCmd.Flags().BoolVar(&Web, "web", false, "view the issue in your browser")
	issueCmd.AddCommand(viewCmd)
}

func PrintSummary(issue *client.Issue) {
	fmt.Println(aurora.Bold(issue.Title))
	var state string
	switch issue.State {
	case "new":
		state = aurora.BgIndex(55, " NEW ").String()
	case "open":
		state = aurora.BgGray(12, " OPEN ").String()
	case "on hold":
		state = aurora.BgBlue(" ON HOLD ").String()
	case "invalid":
		state = aurora.BgRed(" INVALID ").String()
	case "resolved":
		state = aurora.BgGreen(" RESOLVED ").String()
	case "duplicate":
		state = aurora.BgYellow(" DUPLICATE ").String()
	case "wontfix":
		state = aurora.BgRed(" WONTFIX ").String()
	case "closed":
		state = aurora.BgGreen(" CLOSED ").String()
	default:
		state = issue.State
	}

	infoText := aurora.Index(242, fmt.Sprintf("%s opened %s", issue.Repository.FullName, issue.CreatedOn))
	fmt.Printf("%s • %s\n", state, infoText)
	assignee := issue.Assignee.DisplayName
	if assignee == "" {
		assignee = "--"
	}
	fmt.Printf("Type: %s • Priority: %s • Assignee: %s\n", issue.Type, issue.Priority, assignee)

	var thirdLine strings.Builder
	if issue.Component.Name != "" {
		thirdLine.WriteString("Component: ")
		thirdLine.WriteString(issue.Component.Name)
		thirdLine.WriteString(" • ")
	}

	if issue.Milestone.Name != "" {
		thirdLine.WriteString("Milestone: ")
		thirdLine.WriteString(issue.Milestone.Name)
		thirdLine.WriteString(" • ")
	}

	if issue.Version.Name != "" {
		thirdLine.WriteString("Version: ")
		thirdLine.WriteString(issue.Version.Name)
	}

	if thirdLine.Len() > 0 {
		fmt.Println(thirdLine.String())
	}

	if issue.Content.Raw != "" {
		out, err := glamour.Render(issue.Content.Raw, "dark")
		if err != nil {
			logging.Error(err)
			return
		}
		fmt.Println(out)
	}

	footer := aurora.Index(242, fmt.Sprintf("View this issue on Bitbucket.org: %s", issue.Links["html"].Href)).String()
	fmt.Println(footer)
}
