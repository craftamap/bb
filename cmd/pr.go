package cmd

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/cli/cli/git"

	bbgit "github.com/craftamap/bb/git"
)

var (
	prCommand = cobra.Command{
		Use: "pr",
	}
	prListCommand = cobra.Command{
		Use: "list",
		Run: list,
	}
	prViewCommand = cobra.Command{
		Use: "view",
		Run: view,
	}
	prCreateCommand = cobra.Command{
		Use: "create",
		Run: create,
	}

	createOpts = struct {
		Body      string
		Assignees []string
	}{}
)

func init() {
	prCreateCommand.Flags().StringVarP(&createOpts.Body, "body", "b", "", "Supply a body.")
	prCreateCommand.Flags().StringSliceVarP(&createOpts.Assignees, "assignee", "a", nil, "Assign people by their `login`")

	prCommand.AddCommand(&prListCommand)
	prCommand.AddCommand(&prViewCommand)
	prCommand.AddCommand(&prCreateCommand)
	rootCmd.AddCommand(&prCommand)
}

func list(cmd *cobra.Command, args []string) {
	bbrepo, err := bbgit.GetBitbucketRepo()
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}
	prs, err := internal.PrList(globalOpts.Username, globalOpts.Password, bbrepo.RepoOrga, bbrepo.RepoSlug)
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
	}

	fmt.Println()
	fmt.Printf("%sShowing %d of %d open pull requests in %s/%s\n", aurora.Blue(" :: "), len(prs.Values), prs.Size, bbrepo.RepoOrga, bbrepo.RepoSlug)
	fmt.Println()
	for _, pr := range prs.Values {
		fmt.Printf("#%03d  %s   %s -> %s\n", aurora.Green(pr.ID), pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
	}
}

func view(cmd *cobra.Command, args []string) {
	var id int
	var err error

	bbrepo, err := bbgit.GetBitbucketRepo()

	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}

	if len(args) > 0 {
		id, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
			return
		}
	} else {
		branchName, err := git.CurrentBranch()
		if err != nil {
			fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
			return
		}

		prs, err := internal.GetPrIDBySourceBranch(globalOpts.Username, globalOpts.Password, bbrepo.RepoOrga, bbrepo.RepoSlug, branchName)
		if err != nil {
			fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
			return
		}
		if len(prs.Values) == 0 {
			fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Nothing on this branch")
			return
		}

		id = prs.Values[0].ID

	}

	pr, err := internal.PrView(globalOpts.Username, globalOpts.Password, bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}

	fmt.Println(aurora.Bold(pr.Title))
	var state string
	if pr.State == "OPEN" {
		state = aurora.Green("Open").String()
	} else if pr.State == "DECLINED" {
		state = aurora.Red("Declined").String()
	} else {
		state = pr.State
	}

	infoText := aurora.BrightBlack(fmt.Sprintf("%s wants to merge X commits into %s from %s\n", pr.Author.Nickname, pr.Destination.Branch.Name, pr.Source.Branch.Name))
	fmt.Printf("%s • %s\n", state, infoText)
	out, err := glamour.Render(pr.Description, "dark")
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}
	fmt.Println(out)

	footer := aurora.BrightBlack(fmt.Sprintf("View this pull request on Bitbucket.org: %s", pr.Links["html"].Href)).String()
	fmt.Println(footer)
	// fmt.Println(pr, err)

}

func create(cmd *cobra.Command, args []string) {
	bbrepo, err := bbgit.GetBitbucketRepo()
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}

	branchName, err := git.CurrentBranch()
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}

	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}

	fmt.Printf("Creating pull request for %s into %s in %s\n", branchName, "X", fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
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
				Default: branchName,
			},
			Validate: survey.Required,
		},
		{
			Name: "action",
			Prompt: &survey.Select{
				Message: "What's next?",
				Options: []string{"create", "cancel", "continue in browser"},
				Default: "create",
			},
		},
	}
	err = survey.Ask(qs, &answers)
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}

	if answers.Action == "create" {
		response, err := internal.PrCreate(globalOpts.Username, globalOpts.Password, bbrepo.RepoOrga, bbrepo.RepoSlug, branchName, "master", answers.Title)
		if err != nil {
			fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
			return
		}

		fmt.Printf("Take a look at your pull request here:\n")
		fmt.Println(response)
	}
}
