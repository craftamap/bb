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

var BbConfigurationValidation Configuration = map[string]Entry{
	"username": {
		Validator: SimpleStringValidator(),
	},
	"password": {
		Validator: SimpleStringValidator(),
	},
	"remote": {
		Validator: SimpleStringValidator(),
	},
	"git_protocol": {
		Validator: EnumValidator("ssh", "https"),
	},
	"sync-method": {
		Validator: EnumValidator("merge", "rebase"),
	},
}
