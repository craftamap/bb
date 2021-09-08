package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftamap/bb/cmd/options"
	"github.com/craftamap/bb/config"
	"github.com/craftamap/bb/util/logging"
	"github.com/spf13/cobra"
)

var (
	Local bool
	Get   bool
)

func Add(rootCmd *cobra.Command, _ *options.GlobalOptions) {
	configCommand := cobra.Command{
		Use:   "config",
		Short: "configure bb",
		Long:  "configure bb",
		Args: func(cmd *cobra.Command, args []string) error {
			if Get {
				return cobra.ExactArgs(1)(cmd, args)
			} else {
				return cobra.ExactArgs(2)(cmd, args)
			}
		},
		Run: func(_ *cobra.Command, args []string) {
			if Get {
				// TODO: code here
			} else {
				key := args[0]
				inputValue := args[1]

				newValue, err := config.BbConfigurationValidation.ValidateEntry(key, inputValue)

				if err != nil {
					logging.Error(fmt.Sprintf("failed to validate %s: %s", inputValue, err))
					return
				}

				var configDirectory string
				var filename string
				if Local {
					var err error
					configDirectory, filename, err = config.GetLocalConfigurationPath()
					if err != nil {
						logging.Error(err)
						return
					}
				} else {
					configDirectory, filename = config.GetGlobalConfigurationPath()
				}

				// If the directory does not exist, something is off:
				//   - The global configuration directory get's created in root
				//   - The local configuration directory is a repository, which always exists
				if _, err := os.Stat(configDirectory); os.IsNotExist(err) {
					logging.Error(fmt.Sprintf("Expected directory \"%s\", but the directory does not exist", configDirectory))
					return
				}
				path := filepath.Join(configDirectory, filename)
				// If the config itself does not exist, it's fine (although wierd for global) - we create it now
				if _, err := os.Stat(path); os.IsNotExist(err) {
					logging.Note(fmt.Sprintf("Creating config file %s", path))
					fh, err := os.Create(path)
					if err != nil {
						logging.Error(fmt.Sprintf("Unable to create file %s", path))
					}
					fh.Close()
				}

				logging.Debugf("Config file path: %s", path)

				tmpVp, err := config.GetViperForPath(path)
				if err != nil {
					logging.Error(err)
					return
				}

				isSetAlready := tmpVp.IsSet(key)
				oldValue := tmpVp.Get(key)

				if isSetAlready {
					// Don't print old password values
					if strings.ToLower(key) == config.CONFIG_KEY_AUTH_PASSWORD {
						oldValue = "(truncated)"
					}
					logging.Warning(fmt.Sprintf("\"%s\" is already set. This will overwrite the value of \"%s\" from \"%s\" to \"%s\".", key, key, oldValue, newValue))
				}

				logging.Note(fmt.Sprintf("Setting \"%s\" to \"%s\" in %s", key, newValue, path))
				logging.Debugf("%+v", tmpVp.AllSettings())

				// This will most likely save everything as a string
				// TODO: find this out and find a way to save bools and numbers
				tmpVp.Set(key, newValue)
				logging.Debugf("%+v", tmpVp.AllSettings())

				config.WriteViper(tmpVp, path)

				logging.SuccessExclamation(fmt.Sprintf("Successfully updated configuration %s", path))
			}
		},
	}

	configCommand.Flags().BoolVar(&Local, "local", false, "local allows to modify the local configuration")
	configCommand.Flags().BoolVar(&Get, "get", false, "gets a configuration value instead of setting it")

	rootCmd.AddCommand(&configCommand)
}
