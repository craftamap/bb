package update

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/craftamap/bb/cmd/commands/issue/shared"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/util/logging"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	Title       string
	Description string
	Assignee    string
	Kind        string
	Priority    string
	Status      string
	Version     string
	Milestone   string
	Component   string
)

func Add(issueCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Updates a issue",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Initialisation
			var err error
			var id int

			if len(args) > 0 {
				id, err = strconv.Atoi(strings.TrimPrefix(args[0], "#"))
				if err != nil {
					logging.Error(err)
					return
				}
			}

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			io := shared.IssueOptions{}

			currentIssue, err := c.IssuesView(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}

			io.Title = currentIssue.Title
			io.Description = currentIssue.Content.Raw
			io.Assignee = currentIssue.Assignee.UUID
			io.Kind = currentIssue.Kind
			io.Component = currentIssue.Component.Name
			io.Milestone = currentIssue.Milestone.Name
			io.Priority = currentIssue.Priority
			io.Status = currentIssue.State
			io.Version = currentIssue.Version.Name

			// Apply command line args here
			if Title != "" {
				io.Title = Title
			}

			if Description != "" {
				io.Description = Description
			}

			if Kind != "" {
				io.Kind = Kind
			}

			if Priority != "" {
				io.Priority = Priority
			}

			if Status != "" {
				io.Status = Status
			}

			if Version != "" {
				io.Version = Version
			}

			if Milestone != "" {
				io.Milestone = Milestone
			}

			if Component != "" {
				io.Component = Component
			}

			fmt.Printf("Updating issue in %s\n", fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
			fmt.Println()

			io, err, cancel := shared.AskQuestionsForCreateOrUpdate(io, bbrepo, c)
			if err != nil {
				logging.Error(err)
				return
			}
			if cancel {
				return
			}

			response, err := c.IssuesEdit(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id), io)
			if err != nil {
				fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			link, err := url.Parse(response.Repository.Links["html"].(map[string]interface{})["href"].(string))
			link.Path = path.Join(link.Path, "issues", fmt.Sprintf("%d", response.ID))
			if err != nil {
				logging.Error(err)
			}
			fmt.Printf("Take a look at your issue here: %s\n", aurora.Index(242, link.String()))
		},
	}
	updateCmd.Flags().StringVarP(&Description, "description", "b", "", "Supply a description.")
	updateCmd.Flags().StringVarP(&Title, "title", "t", "", "Supply a title.")
	updateCmd.Flags().StringVar(&Kind, "type", "", "Supply a issue type.")
	updateCmd.Flags().StringVar(&Priority, "priority", "", "Supply a issue priority.")
	updateCmd.Flags().StringVar(&Status, "status", "", "Supply a issue status.")
	updateCmd.Flags().StringVar(&Version, "issue-version", "", "Supply a version.")
	updateCmd.Flags().StringVar(&Milestone, "milestone", "", "Supply a milestone.")
	updateCmd.Flags().StringVar(&Component, "component", "", "Supply a component.")
	issueCmd.AddCommand(updateCmd)
}
