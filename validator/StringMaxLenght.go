package validator

import "fmt"

func StringMaxLength(maxLength int) func(string, string) error {
	return func(value string, name string) error {
		if len(value) <= maxLength {
			return nil
		}
		message := fmt.Sprintf("%v must be at most %v characters long", name, maxLength)
		return &ValidationError{message: message}
	}
}
