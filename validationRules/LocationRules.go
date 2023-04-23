package validationRules

import "kando-backend/validator"

func LocationName(v *validator.FluentValidator[string]) {
	v.Add(validator.StringRequired)
}
