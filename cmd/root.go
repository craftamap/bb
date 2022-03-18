package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craftamap/bb/util/logging"

	"github.com/craftamap/bb/client"
	"github.com/craftamap/bb/cmd/commands/api"
	"github.com/craftamap/bb/cmd/commands/auth"
	"github.com/craftamap/bb/cmd/commands/config"
	"github.com/craftamap/bb/cmd/commands/downloads"
	"github.com/craftamap/bb/cmd/commands/issue"
	"github.com/craftamap/bb/cmd/commands/pipelines"
	"github.com/craftamap/bb/cmd/commands/pr"
	"github.com/craftamap/bb/cmd/commands/repo"
	"github.com/craftamap/bb/cmd/options"
	configuration "github.com/craftamap/bb/config"
	bbgit "github.com/craftamap/bb/git"
	"github.com/kirsle/configdir"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version   = ""
	CommitSHA = ""

	rootCmd = &cobra.Command{
		Use:     "bb",
		Short:   "Bitbucket.org CLI",
		Long:    "Work seamlessly with Bitbucket.org from the command line.",
		Example: `$ bb pr list`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			username := viper.GetString(configuration.CONFIG_KEY_AUTH_USERNAME)
			password := viper.GetString(configuration.CONFIG_KEY_AUTH_PASSWORD)
			remoteName := viper.GetString(configuration.CONFIG_KEY_GIT_REMOTE)

			if _, ok := cmd.Annotations["RequiresRepository"]; ok {
				bbrepo, err := bbgit.GetBitbucketRepo(remoteName)
				if err != nil {
					logging.Error(err)
					os.Exit(1)
				}
				if !bbrepo.IsBitbucketOrg() {
					fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), "Are you sure this is a bitbucket repo?")
					os.Exit(1)
				}

				globalOpts.BitbucketRepo = bbrepo
			}

			if _, ok := cmd.Annotations["RequiresClient"]; ok {
				globalOpts.Client = &client.Client{
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
	rootCmd.PersistentFlags().StringVar(&remoteName, "remote", "origin", "if you are in a repository and don't want to interact with the default remote, you can change it")
	rootCmd.PersistentFlags().BoolVar(&logging.PrintDebugLogs, "debug", false, "enabling this flag allows debug logs to be printed")

	err := viper.BindPFlag(configuration.CONFIG_KEY_AUTH_USERNAME, rootCmd.PersistentFlags().Lookup("username"))
	if err != nil {
		logging.Error(err)
		return
	}
	err = viper.BindPFlag(configuration.CONFIG_KEY_AUTH_PASSWORD, rootCmd.PersistentFlags().Lookup("password"))
	if err != nil {
		logging.Error(err)
		return
	}

	err = viper.BindPFlag(configuration.CONFIG_KEY_GIT_REMOTE, rootCmd.PersistentFlags().Lookup("remote"))
	if err != nil {
		logging.Error(err)
		return
	}

	pr.Add(rootCmd, &globalOpts)
	issue.Add(rootCmd, &globalOpts)
	api.Add(rootCmd, &globalOpts)
	downloads.Add(rootCmd, &globalOpts)
	auth.Add(rootCmd, &globalOpts)
	repo.Add(rootCmd, &globalOpts)
	pipelines.Add(rootCmd, &globalOpts)
	config.Add(rootCmd, &globalOpts)

	if CommitSHA != "" {
		vt := rootCmd.VersionTemplate()
		rootCmd.SetVersionTemplate(vt[:len(vt)-1] + " (" + CommitSHA + ")\n")
	}
	if Version == "" {
		Version = "unknown (built from source)"
	}

	rootCmd.Version = Version
}

func initConfig() {
	viper.SetEnvPrefix("bb")
	viper.AutomaticEnv()

	// We support setting the config file manually by running bb with --config.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	} else {
		// create global config directory, first
		configDirectory := configdir.LocalConfig("bb")
		err := configdir.MakePath(configDirectory)
		if err != nil {
			panic(err)
		}
		globalConfigFilePath := filepath.Join(configDirectory, "configuration.toml")
		// create global config directory, first
		if _, err = os.Stat(globalConfigFilePath); os.IsNotExist(err) {
			fh, err := os.Create(globalConfigFilePath)
			if err != nil {
				panic(err)
			}
			defer fh.Close()
		}
		viper.SetConfigType("toml")
		viper.SetConfigName("configuration.toml")
		viper.AddConfigPath(configDirectory)
		err = viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s", err))
		}

		// also read in local configuration
		if repoPath, err := bbgit.RepoPath(); err == nil {
			// the local configuration can be found in the root of a repository
			// If we in a repository, check for the file
			if _, err = os.Stat(filepath.Join(repoPath, ".bb")); err == nil {
				viper.SetConfigType("toml")
				viper.SetConfigName(".bb")
				viper.AddConfigPath(repoPath)
				err = viper.MergeInConfig()
				if err != nil { // Handle errors reading the config file
					panic(fmt.Errorf("fatal error config file: %s", err))
				}
			}
		}
	}

	// Register Aliases for backward compatibility
	viper.RegisterAlias("username", configuration.CONFIG_KEY_AUTH_USERNAME)
	viper.RegisterAlias("password", configuration.CONFIG_KEY_AUTH_PASSWORD)
	viper.RegisterAlias("remote", configuration.CONFIG_KEY_GIT_REMOTE)
	viper.RegisterAlias("git_protocol", configuration.CONFIG_KEY_REPO_CLONE_GIT_PROTOCOL)
	viper.RegisterAlias("sync-method", configuration.CONFIG_KEY_PR_SYNC_SYNC_METHOD)
}
