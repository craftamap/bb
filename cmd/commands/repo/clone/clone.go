package clone

import (
	"fmt"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Add(repoCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	cloneCmd := &cobra.Command{
		Use:   "clone <repoOrga/repoSlug>",
		Short: "clone a repository",
		Annotations: map[string]string{
			"RequiresClient": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := globalOpts.Client

			gitProtocol := viper.GetString("git_protocol")
			if gitProtocol == "" || (gitProtocol != "ssh" && gitProtocol != "https") {
				err := survey.AskOne(&survey.Select{
					Message: "Please select a preferred protocol of cloning repositories",
					Options: []string{"ssh", "https"},
				}, &gitProtocol)
				if err != nil {
					logging.Error(err)
					return
				}
				viper.Set("git_protocol", gitProtocol)
				err = viper.WriteConfig()
				if err != nil {
					logging.Error(err)
					return
				}
			}

			if len(args) == 0 {
				workspaces, err := c.GetWorkspaces()
				if err != nil {
					logging.Error(err)
					return
				}
				// fmt.Println(workspaces)

				workspaceSlugs := []string{}
				for _, workspace := range workspaces.Workspaces {
					workspaceSlugs = append(workspaceSlugs, workspace.Slug)
				}
				workspaceSlug := ""
				err = survey.AskOne(&survey.Select{
					Message: "Select a workspace you want to clone a repository from",
					Options: workspaceSlugs,
				}, &workspaceSlug)
				if err != nil {
					logging.Error(err)
					return
				}
				// fmt.Println(workspaceSlug)
				repos, err := c.RepositoriesForWorkspace(workspaceSlug)
				if err != nil {
					logging.Error(err)
					return
				}
				repoSlugs := []string{}
				for _, repos := range repos {
					repoSlugs = append(repoSlugs, repos.FullName)
				}
				repoOrgaSlug := ""
				err = survey.AskOne(&survey.Select{
					Message: "Select a repository you want to clone",
					Options: repoSlugs,
				}, &repoOrgaSlug)
				if err != nil {
					logging.Error(err)
					return
				}

				splitted := strings.Split(repoOrgaSlug, "/")
				if len(splitted) == 2 {
					_, err := c.RepositoryGet(splitted[0], splitted[1])
					if err != nil {
						logging.Error(err)
						return
					}

					f := FormatRemoteURL(gitProtocol, splitted[0], splitted[1])
					_, err = git.RunClone(f, []string{})
					if err != nil {
						logging.Error(err)
						return
					}
				}
			} else if len(args) == 1 {
				splitted := strings.Split(args[0], "/")
				if len(splitted) == 2 {
					_, err := c.RepositoryGet(splitted[0], splitted[1])
					if err != nil {
						logging.Error(err)
						return
					}

					f := FormatRemoteURL(gitProtocol, splitted[0], splitted[1])
					_, err = git.RunClone(f, []string{})
					if err != nil {
						logging.Error(err)
						return
					}
				} else {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "too less or many args")
					return
				}
			} else {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "too many args")
				return
			}
		},
	}
	repoCmd.AddCommand(cloneCmd)
}

func FormatRemoteURL(protocol string, repoOrga string, repoSlug string) string {
	if protocol == "ssh" {
		return fmt.Sprintf("git@bitbucket.org:%s/%s.git", repoOrga, repoSlug)
	}

	return fmt.Sprintf("https://bitbucket.org/%s/%s.git", repoOrga, repoSlug)
}
