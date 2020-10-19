package cmd

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/craftamap/bb/internal"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
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
)

func init() {
	prCommand.AddCommand(&prListCommand)
	prCommand.AddCommand(&prViewCommand)
	prCommand.AddCommand(&prCreateCommand)
	rootCmd.AddCommand(&prCommand)
}

func list(cmd *cobra.Command, args []string) {
	prs, err := internal.PrList(globalOpts.Username, globalOpts.Password, globalOpts.RepoOrga, globalOpts.RepoSlug)
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
	}

	fmt.Println()
	fmt.Printf("%sShowing %d of %d open pull requests in %s/%s\n", aurora.Blue(" :: "), len(prs.Values), prs.Size, globalOpts.RepoOrga, globalOpts.RepoSlug)
	fmt.Println()
	for _, pr := range prs.Values {
		fmt.Printf("#%03d  %s   %s -> %s\n", aurora.Green(pr.ID), pr.Title, pr.Source.Branch.Name, pr.Destination.Branch.Name)
	}
}

func view(cmd *cobra.Command, args []string) {
	prs, err := internal.GetPrIDBySourceBranch(globalOpts.Username, globalOpts.Password, globalOpts.RepoOrga, globalOpts.RepoSlug, "bb1-branch1")
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
		return
	}
	if len(prs.Values) == 0 {
		fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), "Nothing on this branch")
		return
	}

	pr, err := internal.PrView(globalOpts.Username, globalOpts.Password, globalOpts.RepoOrga, globalOpts.RepoSlug, fmt.Sprintf("%d", prs.Values[0].ID))
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

	footer := aurora.BrightBlack(fmt.Sprintf("View this pull request on Bitbucket.org: %s\n", pr.Links["html"].Href)).String()
	fmt.Printf(footer)
	// fmt.Println(pr, err)

}

func create(cmd *cobra.Command, args []string) {

}
