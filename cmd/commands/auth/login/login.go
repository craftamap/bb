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
	"github.com/spf13/viper"
)

func Add(authCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	loginCmd := &cobra.Command{
		Use: "login",
		Run: func(_ *cobra.Command, _ []string) {
			configDirectory, filename := config.GetGlobalConfigurationPath()
			path := filepath.Join(configDirectory, filename)
			// TODO: extract tmpVp stuff to a seperate file
			tmpVp := viper.New()
			tmpVp.SetConfigType("toml")
			tmpVp.SetConfigFile(path)
			tmpVp.ReadInConfig()

			oldPw := tmpVp.GetString("password")

			if oldPw != "" {
				logging.Warning("You are already logged in as ", tmpVp.GetString("username"))
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
				logging.Error(err)
				return
			}
			username, err := config.BbConfigurationValidation.ValidateEntry("username", answers.Username)
			if err != nil {
				logging.Error(err)
				return
			}
			password, err := config.BbConfigurationValidation.ValidateEntry("password", answers.Password)
			if err != nil {
				logging.Error(err)
				return
			}

			tmpVp.Set("username", username)
			tmpVp.Set("password", password)

			err = tmpVp.WriteConfig()
			if err != nil {
				logging.Error(err)
				return
			}

			logging.Success(fmt.Sprint("Stored credentials successfully to", tmpVp.ConfigFileUsed()))
		},
	}

	authCmd.AddCommand(loginCmd)
}
