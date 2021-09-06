package config

import (
	"fmt"
	"strings"
)

type Validator func(interface{}) (interface{}, error)

// Enum is just a string of a list
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

// SimpleStringValidator validates if a input is a "simple" string - only single-line strings are supported
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

// Entry contains all the data required for Validation and Convertion
type Entry struct {
	Validator Validator
}

type Configuration map[string]Entry

func (c Configuration) ValidateEntry(key string, value interface{}) (interface{}, error) {
	e, ok := c[key]
	if !ok {
		return "", fmt.Errorf("key \"%s\" is not a valid key", key)
	}
	return e.Validator(value)
}

const (
	CONFIG_KEY_AUTH_USERNAME = "auth.username"
	CONFIG_KEY_AUTH_PASSWORD = "auth.password"
	CONFIG_KEY_GIT_REMOTE = "git.remote"
	CONFIG_KEY_REPO_CLONE_GIT_PROTOCOL = "repo.clone.git_protocol"
	CONFIG_KEY_PR_SYNC_SYNC_METHOD = "pr.sync.sync_method"
)


var BbConfigurationValidation Configuration = map[string]Entry{
	CONFIG_KEY_AUTH_USERNAME: {
		Validator: SimpleStringValidator(),
	},
	CONFIG_KEY_AUTH_PASSWORD: {
		Validator: SimpleStringValidator(),
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
