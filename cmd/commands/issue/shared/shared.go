package shared

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
	"github.com/cli/cli/pkg/surveyext"
	"github.com/craftamap/bb/client"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/util/logging"
	"github.com/logrusorgru/aurora"
)

type IssueOptions struct {
	Title       string
	Description string
	Assignee    string
	Kind        string
	Priority    string
	Status      string
	Version     string
	Milestone   string
	Component   string
}

func AskQuestionsForCreateOrUpdate(io IssueOptions, bbrepo *bbgit.BitbucketRepo, c *client.Client) (IssueOptions, error, bool) {
	var err error
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
			return IssueOptions{}, err, false
		}

		if doNext == "cancel" {
			return IssueOptions{}, nil, true
		} else if doNext == "create" {
			break
		}

		if doNext == "modify title" {
			io.Title, err = modifyTitle(io.Title)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "modify description" {
			io.Description, err = modifyDescription(io.Description)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "select assignee" {
			io.Assignee, err = selectAssignee(bbrepo, c, io.Assignee)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "modify kind" {
			io.Kind, err = namedSelectFromOptions("kind", []string{"bug", "enhancement", "proposal", "task"}, io.Kind)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "modify priority" {
			io.Priority, err = namedSelectFromOptions("priority", []string{"trivial", "minor", "major", "critical", "blocker"}, io.Priority)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "modify status" {
			io.Status, err = namedSelectFromOptions("status", []string{"new", "open", "resolved", "invalid", "duplicate", "wontfix", "closed", "on hold"}, io.Status)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "modify component" {
			io.Component, err = namedInput("component", io.Component)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "modify milestone" {
			io.Milestone, err = namedInput("milestone", io.Milestone)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		} else if doNext == "modify version" {
			io.Version, err = namedInput("version", io.Version)
			if err != nil {
				return IssueOptions{}, err, false
			}
			continue
		}
	}

	return io, nil, false
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
		Message:  "Which user do you want to assign?",
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
	logging.SuccessExclamation(fmt.Sprintf("%s: %s", aurora.Bold(name), value))

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
	logging.SuccessExclamation(fmt.Sprintf("%s: %s", aurora.Bold(name), value))

	return value, nil
}
