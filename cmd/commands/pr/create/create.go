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
	Title string
	Body  string
	Force bool
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				sourceBranch string
			)
			// Initialisation
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

			// Get Current Branch
			sourceBranch, err = git.CurrentBranch()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			// Prepare required data
			var (
				targetBranch string
				title        string
				body         string
				defaultBody  string
				reviewers    []string
			)
			// First, init default data
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

			title, body, err = c.PrDefaultTitleAndBody(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch)
			defaultBody = body
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			defaultReviewers, err := c.GetDefaultReviewers(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			for _, rev := range defaultReviewers.Values {
				reviewers = append(reviewers, rev.UUID)
			}

			// Then, check if a pr is already existing. If force is True, take that data
			possiblePrs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			if !Force {
				if len(possiblePrs.Values) != 0 {
					id := possiblePrs.Values[0].ID
					fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), fmt.Sprintf("Pull request %d already exists for this branch. Use --force to ignore this.", id))
					return
				}
			} else {
				if len(possiblePrs.Values) > 0 {
					existingPr := possiblePrs.Values[0]
					title = existingPr.Title
					body = existingPr.Description
					reviewers = []string{}
					for _, reviewer := range existingPr.Reviewers {
						// TODO: make this memory efficient
						reviewers = append(reviewers, reviewer.UUID)
					}
				}
			}
			verb := "Creating"
			if Force {
				verb = "Re-Creating"
			}

			fmt.Printf("%s pull request for %s into %s in %s\n", verb, sourceBranch, targetBranch, fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
			fmt.Println()

			if Title == "" {
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
					err := survey.AskOne(&survey.Input{
						Message: "type your destination branch",
						Default: targetBranch,
					}, &targetBranch)
					if err != nil {
						fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
						return
					}

					// We need to re-generate the body, if the destination branch is changed
					// but only if the body was not modified in before

					_, tempBody, err := c.PrDefaultTitleAndBody(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch)
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
			if err != nil {
				fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			fmt.Printf("Take a look at your pull request here: %s\n", aurora.Index(242, response.Links["html"].Href))
		},
	}
	createCmd.Flags().StringVarP(&Body, "body", "b", "", "Supply a body.")
	createCmd.Flags().StringVarP(&Title, "title", "t", "", "Supply a title.")
	createCmd.Flags().BoolVar(&Force, "force", false, "force creation")
	prCmd.AddCommand(createCmd)
}
