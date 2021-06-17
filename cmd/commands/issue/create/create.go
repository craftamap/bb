package create

import (
	"fmt"
	"net/url"
	"path"

	"github.com/AlecAivazis/survey/v2"
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
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a issue",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Initialisation
			var err error

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			io := shared.IssueOptions{}

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

			fmt.Printf("Creating issue in %s\n", fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
			fmt.Println()

			// If the title was already specified as command line argument, don't ask for it
			if Title == "" {
				questionTitle := &survey.Input{
					Message: "Title",
					Default: io.Title,
				}
				err = survey.AskOne(questionTitle, &io.Title, survey.WithValidator(survey.Required))
			}
			if err != nil {
				logging.Error(err)
				return
			}

			io, err, cancel := shared.AskQuestionsForCreateOrUpdate(io, bbrepo, c)
			if err != nil {
				logging.Error(err)
				return
			}
			if cancel {
				return
			}

			response, err := c.IssuesCreate(bbrepo.RepoOrga, bbrepo.RepoSlug, io)
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
	createCmd.Flags().StringVarP(&Description, "description", "b", "", "Supply a description.")
	createCmd.Flags().StringVarP(&Title, "title", "t", "", "Supply a title.")
	createCmd.Flags().StringVar(&Kind, "type", "", "Supply a issue type.")
	createCmd.Flags().StringVar(&Priority, "priority", "", "Supply a issue priority.")
	createCmd.Flags().StringVar(&Status, "status", "", "Supply a issue status.")
	createCmd.Flags().StringVar(&Version, "issue-version", "", "Supply a version.")
	createCmd.Flags().StringVar(&Milestone, "milestone", "", "Supply a milestone.")
	createCmd.Flags().StringVar(&Component, "component", "", "Supply a component.")
	issueCmd.AddCommand(createCmd)
}
