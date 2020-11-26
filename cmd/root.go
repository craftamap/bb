package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craftamap/bb/cmd/commands/api"
	"github.com/craftamap/bb/cmd/commands/auth"
	"github.com/craftamap/bb/cmd/commands/downloads"
	"github.com/craftamap/bb/cmd/commands/pr"
	"github.com/craftamap/bb/cmd/commands/repo"
	"github.com/craftamap/bb/cmd/options"
	bbgit "github.com/craftamap/bb/git"
	"github.com/craftamap/bb/internal"
	"github.com/kirsle/configdir"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "bb",
		Short:   "Bitbucket.org CLI",
		Long:    "Work seamlessly with Bitbucket.org from the command line.",
		Example: `$ bb pr list`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			username := viper.GetString("username")
			password := viper.GetString("password")

			if _, ok := cmd.Annotations["RequiresRepository"]; ok {
				bbrepo, err := bbgit.GetBitbucketRepo(remoteName)
				if err != nil {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
					os.Exit(1)
				}
				if !bbrepo.IsBitbucketOrg() {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "Are you sure this is a bitbucket repo?")
					os.Exit(1)
				}

				globalOpts.BitbucketRepo = bbrepo
			}

			if _, ok := cmd.Annotations["RequiresClient"]; ok {
				globalOpts.Client = &internal.Client{
					Username: username,
					Password: password,
				}
			}

			if cmd.Name() != "login" {
				if password == "" {
					fmt.Println(aurora.Yellow("::"), aurora.Bold("Warning:"), "Look's like you have not set up bb yet.")
					fmt.Println(aurora.Yellow("::"), aurora.Bold("Warning:"), "Run", aurora.BgWhite(aurora.Black(" bb auth login ")), "to set up bb.")
				}
			}

		},
	}

	cfgFile    string
	globalOpts = options.GlobalOptions{}

	username   string
	password   string
	remoteName string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bb)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "app password")
	rootCmd.PersistentFlags().StringVar(&remoteName, "remote", "origin", "if you are in a repository and don't want to interact with the default origin, you can change it")

	err := viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
		return
	}
	err = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	if err != nil {
		fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), err)
		return
	}

	pr.Add(rootCmd, &globalOpts)
	api.Add(rootCmd, &globalOpts)
	downloads.Add(rootCmd, &globalOpts)
	auth.Add(rootCmd, &globalOpts)
	repo.Add(rootCmd, &globalOpts)
}

func initConfig() {
	if cfgFile == "" {
		configDir := configdir.LocalConfig("bb")
		err := configdir.MakePath(configDir)
		if err != nil {
			panic(err)
		}
		cfgFile = filepath.Join(configDir, "configuration.toml")
		if _, err = os.Stat(cfgFile); os.IsNotExist(err) {
			fh, err := os.Create(cfgFile)
			if err != nil {
				panic(err)
			}
			defer fh.Close()
		}
	}

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
