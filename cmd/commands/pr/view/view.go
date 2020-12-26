package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/client"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	Web bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	viewCmd := &cobra.Command{
		Use:   "view [<nr of pr>]",
		Short: "View a pull request",
		Long:  "Display the title, body, and other information about a pull request.",
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
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
			} else {
				branchName, err := git.CurrentBranch()
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}

				prs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, branchName)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				if len(prs.Values) == 0 {
					fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Nothing on this branch")
					return
				}
				id = prs.Values[0].ID
			}

			pr, err := c.PrView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}
			if Web {
				err := browser.OpenURL(pr.Links["html"].Href)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				return
			}

			commits, err := c.PrCommits(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			PrintSummary(pr, commits)
		},
	}
	viewCmd.Flags().BoolVar(&Web, "web", false, "view the pull request in your browser")
	prCmd.AddCommand(viewCmd)
}

func PrintSummary(pr *client.PullRequest, commits *client.Commits) {
	fmt.Println(aurora.Bold(pr.Title))
	var state string
	if pr.State == "OPEN" {
		state = aurora.Green("Open").String()
	} else if pr.State == "DECLINED" {
		state = aurora.Red("Declined").String()
	} else {
		state = pr.State
	}

	nrOfCommits := len(commits.Values)
	nrOfCommitsStr := fmt.Sprintf("%d", nrOfCommits)
	if nrOfCommits == 10 {
		nrOfCommitsStr = nrOfCommitsStr + "+"
	}

	infoText := aurora.Index(242, fmt.Sprintf("%s wants to merge %s commits into %s from %s\n", pr.Author.Nickname, nrOfCommitsStr, pr.Destination.Branch.Name, pr.Source.Branch.Name))
	fmt.Printf("%s • %s\n", state, infoText)
	if pr.Description != "" {
		out, err := glamour.Render(pr.Description, "dark")
		if err != nil {
			fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
			return
		}
		fmt.Println(out)
	}

	if len(pr.Participants) > 0 {
		fmt.Println("Reviewers: ")
		for _, participant := range pr.Participants {
			if participant.Role == "REVIEWER" {
				var state fmt.Stringer
				switch participant.State {
				case "":
					state = aurora.Index(242, "⌛ PENDING ")
				case "approved":
					state = aurora.Green("✅ APPROVED")
				case "changes_requested":
					state = aurora.Yellow("➖ CHANGE  ")
				}
				fmt.Println(state, participant.User.DisplayName)
			}
		}
	}

	footer := aurora.Index(242, fmt.Sprintf("View this pull request on Bitbucket.org: %s", pr.Links["html"].Href)).String()
	fmt.Println(footer)
}
