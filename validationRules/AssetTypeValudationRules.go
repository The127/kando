package validationRules

import "kando-backend/validator"

func AssetTypeName(v *validator.FluentValidator[string]) {
	v.Add(validator.StringRequired)
}
