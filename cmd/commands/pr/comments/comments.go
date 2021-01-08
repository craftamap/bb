package comments

import (
	"fmt"
	"github.com/craftamap/bb/util/logging"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/cli/cli/git"
	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/client"
	"github.com/kyokomi/emoji"
	"github.com/logrusorgru/aurora"
	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/cobra"
)

var (
	ReviewersNameCache = map[string]string{}
)

func Add(prCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	commentsCmd := &cobra.Command{
		Use:   "comments [<number of pr>]",
		Long:  "View comments",
		Short: "View comments",
		Annotations: map[string]string{
			"RequiresClient":     "true",
			"RequiresRepository": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			var id int
			var err error
			c := globalOpts.Client
			bbrepo := globalOpts.BitbucketRepo
			if len(args) > 0 {
				id, err = strconv.Atoi(strings.TrimPrefix(args[0], "#"))
				if err != nil {
					logging.Error(err)
					return
				}
			} else {
				branchName, err := git.CurrentBranch()
				if err != nil {
					logging.Error(err)
					return
				}

				prs, err := c.GetPrIDBySourceBranch(bbrepo.RepoOrga, bbrepo.RepoSlug, branchName)
				if err != nil {
					logging.Error(err)
					return
				}
				if len(prs.Values) == 0 {
					logging.Warning("Nothing on this branch")
					return
				}

				id = prs.Values[0].ID
			}
			comments, err := c.PrThreadedComments(bbrepo.RepoOrga, bbrepo.RepoSlug, fmt.Sprintf("%d", id))
			if err != nil {
				logging.Error(err)
				return
			}

			members, err := c.GetWorkspaceMembers(bbrepo.RepoOrga)
			if err == nil {
				for _, member := range members.Values {
					ReviewersNameCache[member.AccountID] = member.DisplayName
				}
			}

			for _, comment := range comments {
				ReviewersNameCache[comment.User.AccountID] = comment.User.DisplayName
			}

			re := regexp.MustCompile(`@\{([\w\:-]+)\}`)
			r, _ := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(-1),
			)

			var recursiveComments func(*client.Comment, int)
			recursiveComments = func(comment *client.Comment, depth int) {
				prefix := strings.Repeat(" ", depth) + aurora.Blue("│").String()

				fmt.Println(prefix)
				fmt.Println(prefix, aurora.Blue(aurora.Bold(comment.User.DisplayName)), comment.UpdatedOn)
				if comment.Inline.Path != "" {
					pos := aurora.Index(242, "position: ").String()
					if comment.Inline.From != 0 {
						pos += aurora.Red(fmt.Sprintf("-%d", comment.Inline.From)).String()
					}
					if comment.Inline.To != 0 {
						pos += aurora.Green(fmt.Sprintf("+%d", comment.Inline.To)).String()
					}
					fmt.Println(prefix, comment.Inline.Path, pos)
				}

				commentContent := comment.Content.Raw
				commentContent = emoji.Sprint(commentContent)
				commentContent, _ = r.Render(commentContent)
				occurrences := re.FindAllStringSubmatch(commentContent, -1)
				for _, occ := range occurrences {
					name, ok := ReviewersNameCache[occ[1]]
					if ok {
						commentContent = strings.Replace(commentContent, occ[0], aurora.Bold("@"+name).String(), -1)
					}
				}
				commentContent = wordwrap.String(commentContent, 120-depth)
				for _, line := range strings.Split(commentContent, "\n") {
					fmt.Println(prefix, line)
				}

				sort.SliceStable(comment.Children, func(i int, j int) bool {
					return strings.Compare(comment.Children[i].CreatedOn, comment.Children[j].CreatedOn) < 0
				})

				for _, comment := range comment.Children {
					recursiveComments(comment, depth+1)
				}
			}

			sort.SliceStable(comments, func(i int, j int) bool {
				return strings.Compare(comments[i].CreatedOn, comments[j].CreatedOn) < 0
			})
			for _, comment := range comments {
				recursiveComments(comment, 0)
				fmt.Println()
			}
		},
	}
	prCmd.AddCommand(commentsCmd)
}
