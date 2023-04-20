package validator

import "fmt"

func Required(value any, name string) error {
	if value != "" {
		return nil
	}
	message := fmt.Sprintf("%v is required", name)
	return &ValidationError{message: message}
}
