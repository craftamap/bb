package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/craftamap/bb/util/logging"
	"github.com/spf13/viper"
)

const (
	CONFIG_KEY_AUTH_USERNAME           = "auth.username"
	CONFIG_KEY_AUTH_PASSWORD           = "auth.password"
	CONFIG_KEY_GIT_REMOTE              = "git.remote"
	CONFIG_KEY_REPO_CLONE_GIT_PROTOCOL = "repo.clone.git_protocol"
	CONFIG_KEY_PR_SYNC_SYNC_METHOD     = "pr.sync.sync_method"
)

type Validator func(interface{}) (interface{}, error)

// Enum is just a string of a list.
func EnumValidator(validValues ...string) Validator {
	return func(inputValue interface{}) (interface{}, error) {
		_, ok := inputValue.(string)
		if !ok {
			return "", fmt.Errorf("value \"%s\" is not a string, but of type %T", inputValue, inputValue)
		}
		isInList := false
		for _, validValue := range validValues {
			if inputValue == validValue {
				isInList = true
				break
			}
		}
		if !isInList {
			return "", fmt.Errorf("value \"%s\" is not a valid value. Valid Values are %s", inputValue, validValues)
		}

		return inputValue, nil
	}
}

// SimpleStringValidator validates if a input is a "simple" string - only single-line strings are supported.
func SimpleStringValidator() Validator {
	return func(inputValue interface{}) (interface{}, error) {
		_, ok := inputValue.(string)
		if !ok {
			return "", fmt.Errorf("value \"%s\" is not a string, but of type %T", inputValue, inputValue)
		}

		if strings.ContainsAny(inputValue.(string), "\r\n") {
			return "", fmt.Errorf("value \"%s\" contains illegal line break", inputValue)
		}

		return inputValue, nil
	}
}

// Entry contains all the data required for Validation and Convertion.
type Entry struct {
	Validator Validator
	Hidden    bool
}

type Configuration map[string]Entry

var BbConfigurationValidation Configuration = map[string]Entry{
	CONFIG_KEY_AUTH_USERNAME: {
		Validator: SimpleStringValidator(),
	},
	CONFIG_KEY_AUTH_PASSWORD: {
		Validator: SimpleStringValidator(),
		Hidden:    true,
	},
	CONFIG_KEY_GIT_REMOTE: {
		Validator: SimpleStringValidator(),
	},
	CONFIG_KEY_REPO_CLONE_GIT_PROTOCOL: {
		Validator: EnumValidator("ssh", "https"),
	},
	CONFIG_KEY_PR_SYNC_SYNC_METHOD: {
		Validator: EnumValidator("merge", "rebase"),
	},
}

func (c Configuration) ValidateEntry(key string, value interface{}) (interface{}, error) {
	e, ok := c[key]
	if !ok {
		return "", fmt.Errorf("key \"%s\" is not a valid key", key)
	}
	return e.Validator(value)
}

func ValidateEntry(key string, value interface{}) (interface{}, error) {
	return BbConfigurationValidation.ValidateEntry(key, value)
}

func ValidateAndUpdateEntry(filepath string, key string, value interface{}) (interface{}, error) {
	sanitizedValue, err := ValidateEntry(key, value)
	if err != nil {
		return "", err
	}

	// TODO: Add a filename-to-tmpVp cache - this way, we can prevent creating a new viper every time we want to set a value
	vp, err := GetViperForPath(filepath)
	if err != nil {
		return sanitizedValue, err
	}

	vp.Set(key, sanitizedValue)
	err = WriteViper(vp, filepath)

	return sanitizedValue, err
}

func ValidateAndUpdateEntryWithViper(vp viper.Viper, key string, value interface{}) (interface{}, error) {
	sanitizedValue, err := ValidateEntry(key, value)
	if err != nil {
		return "", err
	}

	vp.Set(key, sanitizedValue)
	err = WriteViper(vp, vp.ConfigFileUsed())

	return sanitizedValue, err
}

func WriteViper(vp viper.Viper, path string) error {
	// WORKAROUND: currently, WriteConfig does not support writing to `.bb`-files despite setting SetConfigType.
	// Therefore, we create a temporary file, write there, and try to copy the file over.
	tmpFh, err := ioutil.TempFile(os.TempDir(), "bb-tmpconfig.*.toml")
	if err != nil {
		logging.Error("Failed to create temporary configuration file")
		return err
	}
	tmpFilename := tmpFh.Name()
	logging.Debugf("tmpFilename: %s", tmpFilename)
	err = tmpFh.Close()
	if err != nil {
		logging.Error("Failed to create temporary configuration file")
		return err
	}
	err = vp.WriteConfigAs(tmpFilename)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to write temporary config %s: %s", path, err))
		return err
	}
	err = copyFileContent(tmpFilename, path)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to write config %s -> %s: %s", tmpFilename, path, err))
		return err
	}
	return nil
}

func copyFileContent(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst) // Create or trunicate
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func GetViperForPath(path string) (viper.Viper, error) {
	tmpVp := viper.New()
	tmpVp.SetConfigType("toml")
	tmpVp.SetConfigFile(path)
	err := tmpVp.ReadInConfig()

	return *tmpVp, err
}
