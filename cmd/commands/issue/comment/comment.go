package comment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/glamour"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/util/editor"
	"github.com/craftamap/bb/util/logging"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func Add(issueCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	commentCmd := &cobra.Command{
		Use:   "comment [<nr of issue>]",
		Short: "comment a issue",
		Long:  "Add a comment to an issue",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error

			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo

			if len(args) == 1 {
				id, err = strconv.Atoi(strings.TrimPrefix(args[0], "#"))
				if err != nil {
					logging.Error(err)
					return
				}
			} else {
				logging.Error("Wrong number of arguments, see --help")
				return
			}

			body, err := editor.OpenInEditor("", "bb-issuecomment*.md")
			if err != nil {
				logging.Error(err)
				return
			}

			fmt.Println(aurora.Bold(aurora.Green("!").String() + " Body:"))

			out, _ := glamour.Render(body, "dark")
			fmt.Print(out)

			var confirmation bool
			err = survey.AskOne(&survey.Confirm{
				Message: "Do you want to comment this?",
				Default: true,
			}, &confirmation)
			if err != nil {
				logging.Error(err)
				return
			}

			if !confirmation {
				return
			}

			response, err := c.IssuesComment(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id), body)
			if err != nil {
				logging.Error(err)
				return
			}

			fmt.Printf("Take a look at your comment here: %s\n", aurora.Index(242, response.Links["html"].Href))
		},
	}
	issueCmd.AddCommand(commentCmd)
}
