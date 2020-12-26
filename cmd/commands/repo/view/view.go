package view

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/client"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	Web bool
)

func Add(repoCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View a pull request",
		Long:  "Display the title, body, and other information about a pull request.",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},

		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			repo, err := c.RepositoryGet(bbrepo.RepoOrga, bbrepo.RepoSlug)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}
			if Web {
				err := browser.OpenURL(repo.Links["html"].(map[string]interface{})["href"].(string))
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}
				return
			}

			readme, err := c.GetReadmeContent(bbrepo.RepoOrga, bbrepo.RepoSlug, repo.MainBranch.Name)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			PrintSummary(repo, readme)
		},
	}

	viewCmd.Flags().BoolVar(&Web, "web", false, "view this repository in your web browser")
	repoCmd.AddCommand(viewCmd)
}

func PrintSummary(repo *client.Repository, readme string) {
	fmt.Println(aurora.Bold(repo.FullName))
	if repo.Description != "" {
		fmt.Println(aurora.Bold("Description:"), repo.Description)
	}
	fmt.Println(aurora.Bold("Project:"), repo.Project.Name, fmt.Sprintf("(%s)", repo.Project.Key))
	if readme != "" {
		out, err := glamour.Render(readme, "dark")
		if err != nil {
			fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
			return
		}
		fmt.Println(out)
	}
	footer := aurora.Index(242, fmt.Sprintf("View this repository on Bitbucket.org: %s", repo.Links["html"].(map[string]interface{})["href"].(string))).String()
	fmt.Println(footer)
}
