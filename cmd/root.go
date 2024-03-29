package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/craftamap/bb/client"
	"github.com/craftamap/bb/cmd/commands/api"
	"github.com/craftamap/bb/cmd/commands/auth"
	"github.com/craftamap/bb/cmd/commands/downloads"
	"github.com/craftamap/bb/cmd/commands/issue"
	"github.com/craftamap/bb/cmd/commands/pipelines"
	"github.com/craftamap/bb/cmd/commands/pr"
	"github.com/craftamap/bb/cmd/commands/repo"
	"github.com/craftamap/bb/cmd/options"
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
			if _, ok := cmd.Annotations["RequiresRepository"]; ok {
				if repository != "" {
					if _, ok := cmd.Annotations["RequiresFSRepository"]; ok {
						logging.Error("this command requires to be run in the repository directory directly.")
						os.Exit(1)
					}

					splitted := strings.Split(repository, "/")
					if len(splitted) != 2 {
						logging.Error(fmt.Sprintf("repository %s is not in valid format. Please pass it in the repoOrga/repoSlug format.", repository))
						os.Exit(1)
					}
					globalOpts.BitbucketRepo = &bbgit.BitbucketRepo{
						RepoOrga: splitted[0],
						RepoSlug: splitted[1],
					}
					globalOpts.IsFSRepo = false
				} else {
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
					globalOpts.IsFSRepo = true
				}

			}

			username := viper.GetString("username")
			password := viper.GetString("password")
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
	repository string
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
	rootCmd.PersistentFlags().StringVar(&repository, "repository", "", "for most commands, the repository can be specified in the repoOrga/repoSlug format")
	rootCmd.PersistentFlags().StringVar(&remoteName, "remote", "origin", "if you are in a repository and don't want to interact with the default remote, you can change it")
	rootCmd.PersistentFlags().BoolVar(&logging.PrintDebugLogs, "debug", false, "enabling this flag allows debug logs to be printed")

	err := viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	if err != nil {
		logging.Error(err)
		return
	}
	err = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
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

	viper.SetEnvPrefix("bb")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
