package validator

import (
	"fmt"
	"strings"
)

func Email(value string, name string) error {
	if strings.Contains(value, "@") {
		return nil
	}
	message := fmt.Sprintf("%v must be a valid email address", name)
	return &ValidationError{message: message}
}
