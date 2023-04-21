package validationRules

import "kando-backend/validator"

func ManufacturerName(v *validator.FluentValidator[string]) {
	v.Add(validator.StringRequired)
}
