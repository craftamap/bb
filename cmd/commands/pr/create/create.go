package create

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
	"github.com/cli/cli/git"
	"github.com/cli/cli/pkg/surveyext"
	"github.com/cli/cli/utils"
	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	Title       string
	Body        string
	Destination string
	Force       bool
)

var (
	ReviewersNameCache = map[string]string{}
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Initialisation

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			var (
				sourceBranch string
				targetBranch string
				title        string
				body         string
				defaultBody  string
				reviewers    []string
			)
			sourceBranch, err := git.CurrentBranch()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			// Prepare required data
			// First, init default data
			repo, err := c.RepositoryGet(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}
			targetBranch = repo.MainBranch.Name

			if Destination != "" {
				targetBranch = Destination
			}

			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			if _, err := c.GetBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch); err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "This current branch is not available on bitbucket.org yet. You need to push the branch, first.")
				return
			}

			head, err := bbgit.CurrentHead()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			if _, err := c.GetCommit(bbrepo.RepoOrga, bbrepo.RepoSlug, head); err != nil {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Current commit is not available on bitbucket yet. If you create the pull request now, it won't contain the latest pushes.")
			}

			if ucc, err := git.UncommittedChangeCount(); err == nil && ucc > 0 {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), utils.Pluralize(ucc, "uncommitted change"))
			}

			title, body, err = c.PrDefaultTitleAndBody(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch)
			defaultBody = body
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			defaultReviewers, err := c.GetDefaultReviewers(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			currentUser, err := c.GetCurrentUser()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Can't get the current user - this means that the default reviewers won't be added to this pull request. Make sure to grant the account-scope for your access token")
			} else {
				for _, rev := range defaultReviewers.Values {
					if currentUser.Uuid != rev.UUID {
						reviewers = append(reviewers, rev.UUID)
					}
					// Add the user to the cache in any case
					ReviewersNameCache[rev.UUID] = rev.DisplayName
				}
			}

			// Then, check if a pr is already existing. If force is True, take that data
			possiblePrs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
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
					// Only apply the old title if one of the command line options is not set
					if Title == "" && Body == "" && Destination == "" {
						title = existingPr.Title
						body = existingPr.Description
						targetBranch = existingPr.Destination.Branch.Name
					}
					reviewers = []string{}
					for _, reviewer := range existingPr.Reviewers {
						// TODO: make this memory efficient
						reviewers = append(reviewers, reviewer.UUID)
					}
				}
			}

			// Apply command line args here
			if Title != "" {
				title = Title
			}

			if Body != "" {
				body = Body
			}

			verb := "Creating"
			if Force {
				verb = "Re-Creating"
			}

			fmt.Printf("%s pull request for %s into %s in %s\n", verb, sourceBranch, targetBranch, fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
			fmt.Println()

			// If the title was already specified as command line argument, don't ask for it
			if Title == "" {
				questionTitle := &survey.Input{
					Message: "Title",
					Default: title,
				}
				err = survey.AskOne(questionTitle, &title)
			}
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			fmt.Println(aurora.Bold(aurora.Green("!").String() + " Body:"))

			out, _ := glamour.Render(body, "dark")
			fmt.Print(out)

			if len(reviewers) > 0 {
				fmt.Println("Reviewers:")
				for _, reviewer := range reviewers {
					name, ok := ReviewersNameCache[reviewer]
					if ok {
						fmt.Println("-", name)
					} else {
						fmt.Println("-", reviewer)
					}
				}
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
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
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
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
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
						fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}

					// We need to re-generate the body, if the destination branch is changed
					// but only if the body was not modified in before

					_, tempBody, err := c.PrDefaultTitleAndBody(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch)
					if err != nil {
						fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					if body == defaultBody {
						body = tempBody
					}
				}
			}

			response, err := c.PrCreate(bbrepo.RepoOrga, bbrepo.RepoSlug, sourceBranch, targetBranch, title, body, reviewers)
			if err != nil {
				fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			fmt.Printf("Take a look at your pull request here: %s\n", aurora.Index(242, response.Links["html"].Href))
		},
	}
	createCmd.Flags().StringVarP(&Body, "body", "b", "", "Supply a body.")
	createCmd.Flags().StringVarP(&Title, "title", "t", "", "Supply a title.")
	createCmd.Flags().StringVarP(&Destination, "destination", "d", "", "Supply the destination branch of your pull request. Defaults to default branch of the repository")
	createCmd.Flags().BoolVar(&Force, "force", false, "force creation")
	prCmd.AddCommand(createCmd)
}
