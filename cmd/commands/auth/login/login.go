package login

import (
	"fmt"
	"path/filepath"

	"github.com/craftamap/bb/config"
	"github.com/craftamap/bb/util/logging"

	"github.com/AlecAivazis/survey/v2"
	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func Add(authCmd *cobra.Command, _ *options.GlobalOptions) {
	loginCmd := &cobra.Command{
		Use: "login",
		Run: func(_ *cobra.Command, _ []string) {
			configDirectory, filename := config.GetGlobalConfigurationPath()
			path := filepath.Join(configDirectory, filename)
			tmpVp, err := config.GetViperForPath(path)
			if err != nil {
				logging.Error(err)
				return
			}

			oldPw := tmpVp.GetString(config.CONFIG_KEY_AUTH_PASSWORD)

			if oldPw != "" {
				logging.Warning("You are already logged in as ", tmpVp.GetString(config.CONFIG_KEY_AUTH_USERNAME))
				cont := false
				err := survey.AskOne(&survey.Confirm{Message: "Do you want to overwrite this?"}, &cont)
				if err != nil {
					logging.Error(err)
					return
				}

				if !cont {
					return
				}
			}

			logging.Success("In order to use bb, you need to create an app password for bitbucket.org. Navigate to")
			logging.Success(aurora.Index(242, "https://bitbucket.org/account/settings/app-passwords/"))
			logging.Success("And create an app password for your account with the required permissions.")

			answers := struct {
				Username string
				Password string
			}{}

			err = survey.Ask([]*survey.Question{
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
				logging.Error(err)
				return
			}
			_, err = config.ValidateAndUpdateEntryWithViper(tmpVp, config.CONFIG_KEY_AUTH_USERNAME, answers.Username)
			if err != nil {
				logging.Error(err)
				return
			}
			_, err = config.ValidateAndUpdateEntryWithViper(tmpVp, config.CONFIG_KEY_AUTH_PASSWORD, answers.Password)
			if err != nil {
				logging.Error(err)
				return
			}

			logging.Success(fmt.Sprint("Stored credentials successfully to", tmpVp.ConfigFileUsed()))
		},
	}

	authCmd.AddCommand(loginCmd)
}
