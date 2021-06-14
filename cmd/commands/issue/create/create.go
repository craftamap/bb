package create

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
	"github.com/cli/cli/pkg/surveyext"
	"github.com/craftamap/bb/client"
	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
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

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			var (
				title       string
				description string
				assignee    string
				err         error
				kind        string
				priority    string
				status      string
				version     string
				milestone   string
				component   string
			)
			// Apply command line args here
			if Title != "" {
				title = Title
			}

			if Description != "" {
				description = Description
			}

			if Kind != "" {
				kind = Kind
			}

			if Priority != "" {
				priority = Priority
			}

			if Status != "" {
				status = Status
			}

			if Version != "" {
				version = Version
			}

			if Milestone != "" {
				milestone = Milestone
			}

			if Component != "" {
				component = Component
			}

			fmt.Printf("Creating issue in %s\n", fmt.Sprintf("%s/%s", bbrepo.RepoOrga, bbrepo.RepoSlug))
			fmt.Println()

			// If the title was already specified as command line argument, don't ask for it
			if Title == "" {
				questionTitle := &survey.Input{
					Message: "Title",
					Default: title,
				}
				err = survey.AskOne(questionTitle, &title, survey.WithValidator(survey.Required))
			}
			if err != nil {
				logging.Error(err)
				return
			}

			for {
				selectNext := &survey.Select{
					Message: "What's next?",
					Options: []string{
						"create",
						"modify title",
						"modify description",
						"select assignee",
						"modify kind",
						"modify priority",
						"modify status",
						"modify component",
						"modify milestone",
						"modify version",
						"cancel",
					},
					Default: "create",
				}
				var doNext string
				err = survey.AskOne(selectNext, &doNext)
				if err != nil {
					logging.Error(err)
					return
				}

				if doNext == "cancel" {
					return
				} else if doNext == "create" {
					break
				}

				if doNext == "modify title" {
					title, err = modifyTitle(title)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "modify description" {
					description, err = modifyDescription(description)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "select assignee" {
					assignee, err = selectAssignee(bbrepo, c, assignee)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "modify kind" {
					kind, err = namedSelectFromOptions("kind", []string{"bug", "enhancement", "proposal", "task"}, kind)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "modify priority" {
					priority, err = namedSelectFromOptions("priority", []string{"trivial", "minor", "major", "critical", "blocker"}, priority)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "modify status" {
					status, err = namedSelectFromOptions("status", []string{"new", "open", "resolved", "invalid", "duplicate", "wontfix", "closed", "on hold"}, status)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "modify component" {
					component, err = namedInput("component", component)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "modify milestone" {
					milestone, err = namedInput("milestone", milestone)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				} else if doNext == "modify version" {
					version, err = namedInput("version", version)
					if err != nil {
						fmt.Printf("%s%s%#v\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
						return
					}
					continue
				}
			}

			response, err := c.IssuesCreate(bbrepo.RepoOrga, bbrepo.RepoSlug, struct {
				Title       string
				Description string
				Assignee    string
				Kind        string
				Priority    string
				Status      string
				Version     string
				Milestone   string
				Component   string
			}{
				Title:       title,
				Description: description,
				Assignee:    assignee,
				Kind:        kind,
				Priority:    priority,
				Status:      status,
				Version:     version,
				Milestone:   milestone,
				Component:   component,
			})
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

func modifyDescription(body string) (string, error) {
	body, err := surveyext.Edit("vim", "bb-issue*.md", body, os.Stdin, os.Stdout, os.Stderr, nil)
	if err != nil {
		return "", err
	}

	fmt.Println(aurora.Bold(aurora.Green("!").String() + " Body:"))

	out, _ := glamour.Render(body, "dark")
	fmt.Print(out)
	return body, nil
}

func modifyTitle(title string) (string, error) {
	questionTitle := &survey.Input{
		Message: "Title",
		Default: title,
	}
	err := survey.AskOne(questionTitle, &title, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}
	fmt.Println(aurora.Bold(aurora.Green("!").String()+" Title:"), title)
	return title, nil
}

func selectAssignee(bbrepo *bbgit.BitbucketRepo, c *client.Client, assignee string) (string, error) {
	logging.Note("Currently, only members of the current workspace can be added as reviewers.")
	logging.Note("Currently, there is no way of detecting if a user of your workspace has access to the repository. Adding a wrong user without access to the repository leads to a error while creating the repository.")

	members, err := c.GetWorkspaceMembers(bbrepo.RepoOrga)
	if err != nil {
		logging.Warning(fmt.Sprint("Could not get workspace members - create the issue without a assignee and assign them manually using the browser", err))
		return assignee, nil
	}
	logging.Debugf("members: %+v", members)
	nameToUUID := map[string]string{}
	listOfNames := make([]string, 0, len(members.Values))
	assigneeName := ""

	for _, member := range members.Values {
		listOfNames = append(listOfNames, member.User.DisplayName)
		nameToUUID[member.User.DisplayName] = member.User.UUID

		if assignee == member.User.UUID {
			assigneeName = member.User.DisplayName
		}
	}

	listOfNames = append(listOfNames, "(No Assignee)")
	nameToUUID["(No Assignee)"] = ""

	err = survey.AskOne(&survey.Select{
		Message:  "Which user do you wnat to assign?",
		Options:  listOfNames,
		PageSize: 20,
	}, &assigneeName)
	if err != nil {
		logging.Error(err)
		return assignee, err
	}
	assignee = nameToUUID[assigneeName]
	return assignee, nil
}

func namedSelectFromOptions(name string, options []string, value string) (string, error) {
	question := &survey.Select{
		Message: name,
		Default: value,
		Options: options,
	}
	err := survey.AskOne(question, &value)
	if err != nil {
		return "", err
	}
	logging.Success(name, value)

	return value, nil
}

func namedInput(name string, value string) (string, error) {
	question := &survey.Input{
		Message: name,
		Default: value,
	}
	err := survey.AskOne(question, &value)
	if err != nil {
		return "", err
	}
	logging.Success(name, value)

	return value, nil
}
