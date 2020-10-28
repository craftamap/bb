package create

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
	"github.com/cli/cli/git"
	"github.com/cli/cli/pkg/surveyext"
	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	Body string
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	createCmd := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				sourceBranch string
				targetBranch string
				title        string
				defaultBody  string
				body         string
				reviewers    []string
			)

			c := internal.Client{
				Username: globalOpts.Username,
				Password: globalOpts.Password,
			}

			bbrepo, err := bbgit.GetBitbucketRepo()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			if !bbrepo.IsBitbucketOrg() {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Are you sure this is a bitbucket repo?")
				return
			}

			sourceBranch, err = git.CurrentBranch()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			repo, err := c.RepositoryGet(bbrepo.RepoOrga, bbrepo.RepoSlug)
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
					Default: sourceBranch,
				}
				err = survey.AskOne(questionTitle, &title)
			}
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			if Body == "" {
				defaultBody, err = c.PrDefaultBody(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
					return
				}
				body = defaultBody
				fmt.Println(body)
			} else {
				body = Body
			}

			fmt.Println(aurora.Bold(aurora.Green("!").String() + " Body:"))

			out, _ := glamour.Render(body, "dark")
			fmt.Print(out)

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
					body, err = surveyext.Edit("vim", "", body, os.Stdin, os.Stdout, os.Stderr, nil)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
						return
					}

					continue
				}

				if doNext == "change destination branch" {
					survey.AskOne(&survey.Input{
						Message: "type your destination branch",
						Default: targetBranch,
					}, &targetBranch)

					// We need to re-generate the body, if the destination branch is changed
					// but only if the body was not modified in before

					tempBody, err := c.PrDefaultBody(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch)
					if err != nil {
						fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
						return
					}
					if body == defaultBody {
						body = tempBody
					}

				}

			}

			response, err := c.PrCreate(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch, title, body, reviewers)
			fmt.Println(targetBranch)
			if err != nil {
				fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			fmt.Printf("Take a look at your pull request here:\n")
			fmt.Println(response.Links["html"].Href)

		},
	}
	createCmd.Flags().StringVarP(&Body, "body", "b", "", "Supply a body.")
	prCmd.AddCommand(createCmd)
}
