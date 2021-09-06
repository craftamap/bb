package config

import (
	bbgit "github.com/craftamap/bb/git"
	"github.com/kirsle/configdir"
)

// TODO: Extract to util
func GetGlobalConfigurationPath() (configDirectory string, filename string) {
	configDirectory = configdir.LocalConfig("bb")
	return configDirectory, "configuration.toml"
}

func GetLocalConfigurationPath() (configDirectory, filename string, err error) {
	configDirectory, err = bbgit.RepoPath()
	return configDirectory, ".bb", err
}
