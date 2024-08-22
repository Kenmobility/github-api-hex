package helpers

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// IsLocal returns true or false depending on APP_ENV environmental variable's value
func IsLocal() bool {
	return os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "local"
}

// Getenv gets the env variable value or set a default if empty
func Getenv(variable string, defaultValue ...string) string {
	env := os.Getenv(variable)
	if env == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return env
}

func IsRepositoryNameValid(repoName string) bool {
	return strings.Contains(repoName, "/")
}

func ValidateInput(input interface{}) []string {
	var errors []string
	v := validator.New()

	err := v.Struct(input)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			switch e.ActualTag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s field is required", e.Field()))
			default:
				errors = append(errors, "an error occurred")
			}
		}
	}

	return errors
}
