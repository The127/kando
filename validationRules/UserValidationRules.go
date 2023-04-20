package validationRules

import "kando-backend/validator"

func UserDisplayName(v *validator.FluentValidator[string]) {
	v.Add(validator.StringMaxLength(100))
}

func UserUsername(v *validator.FluentValidator[string]) {
	v.Add(validator.StringMaxLength(100))
}
