package validator

import "fmt"

func StringRequired(value string, name string) error {
	if value != "" {
		return nil
	}
	message := fmt.Sprintf("%v is required", name)
	return &ValidationError{message: message}
}
