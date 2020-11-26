package login

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Add(authCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	loginCmd := &cobra.Command{
		Use: "login",
		Run: func(cmd *cobra.Command, args []string) {
			oldPw := viper.GetString("password")

			if oldPw != "" {
				fmt.Println(aurora.Yellow("::"), aurora.Bold("Warning:"), "You are already logged in as", viper.GetString("username"))
				cont := false
				err := survey.AskOne(&survey.Confirm{Message: "Do you want to overwrite this?"}, &cont)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					return
				}

				if !cont {
					return
				}
			}

			fmt.Println(aurora.Green("::"), "In order to use bb, you need to create an app password for bitbucket.org. Navigate to")
			fmt.Println(aurora.Green("::"), aurora.Index(242, "https://bitbucket.org/account/settings/app-passwords/"))
			fmt.Println(aurora.Green("::"), "And create an app password for your account with the required permissions.")

			answers := struct {
				Username string
				Password string
			}{}

			err := survey.Ask([]*survey.Question{
				{
					Name: "username",
					Prompt: &survey.Input{
						Message: "Please enter your username:",
					},
				},
				{
					Name: "password",
					Prompt: &survey.Password{
						Message: "Please enter the app password you just created:",
					},
				},
			}, &answers)

			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			viper.Set("username", answers.Username)
			viper.Set("password", answers.Password)

			err = viper.WriteConfig()
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
				return
			}

			fmt.Println(aurora.Green("::"), "Stored credentials successfully to", viper.ConfigFileUsed())
		},
	}

	authCmd.AddCommand(loginCmd)
}
