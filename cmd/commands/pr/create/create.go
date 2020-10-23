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

			fmt.Printf("Creating pull request for %s into %s in %s\n", sourceBranch, targetBranch, fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
			fmt.Println()

			if title == "" {
				questionTitle := &survey.Input{
					Message: "Title",
					Default: title,
				}
				err = survey.AskOne(questionTitle, &title)
			}
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			for {
				selectNext := &survey.Select{
					Message: "What's next?",
					Options: []string{"create", "modify body", "change destination branch", "cancel"},
					Default: "create",
				}
				var doNext string
				err = survey.AskOne(selectNext, &doNext)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}

				if doNext == "cancel" {
					return
				} else if doNext == "create" {
					break
				}

				if doNext == "modify body" {
					continue
				}

				if doNext == "change destination branch" {
					survey.AskOne(&survey.Input{
						Message: "type your destination branch",
						Default: targetBranch,
					}, &targetBranch)
				}

			}

			response, err := internal.PrCreate(globalOpts.Username, globalOpts.Password, bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch, title, body, reviewers)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			fmt.Printf("Take a look at your pull request here:\n")
			fmt.Println(response)

		},
	}
	createCmd.Flags().StringVarP(&Body, "body", "b", "", "Supply a body.")
	createCmd.Flags().StringSliceVarP(&Assignees, "assignee", "a", nil, "Assign people by their `login`")
	prCmd.AddCommand(createCmd)
}
