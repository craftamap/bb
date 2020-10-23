package create

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	Body      string
	Assignees []string
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	createCmd := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				sourceBranch string
				targetBranch string
				title        string
				body         string
				reviewers    []string
			)

			bbrepo, err := bbgit.GetBitbucketRepo()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			sourceBranch, err = git.CurrentBranch()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			repo, err := internal.RepositoryGet(globalOpts.Username, globalOpts.Password, bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			targetBranch = repo.MainBranch.Name

			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			fmt.Printf("Creating pull request for %s into %s in %s\n", sourceBranch, "X", fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
			fmt.Println()

			answers := struct {
				Title  string
				Action string
			}{}

			// body := createOpts.Body

			var qs = []*survey.Question{
				{
					Name: "title",
					Prompt: &survey.Input{
						Message: "Title",
						Default: sourceBranch,
					},
					Validate: survey.Required,
				},
				{
					Name: "action",
					Prompt: &survey.Select{
						Message: "What's next?",
						Options: []string{"create", "modify body", "set reviewers", "change destination branch", "cancel"},
						Default: "create",
					},
				},
			}
			err = survey.Ask(qs, &answers)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			title = answers.Title

			if answers.Action == "create" {
				response, err := internal.PrCreate(globalOpts.Username, globalOpts.Password, bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch, title, body, reviewers)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}

				fmt.Printf("Take a look at your pull request here:\n")
				fmt.Println(response)
			}

		},
	}
	createCmd.Flags().StringVarP(&Body, "body", "b", "", "Supply a body.")
	createCmd.Flags().StringSliceVarP(&Assignees, "assignee", "a", nil, "Assign people by their `login`")
	prCmd.AddCommand(createCmd)
}
