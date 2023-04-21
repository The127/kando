package fake

type FieldValues struct {
	Values map[string]any
}

func WithDefaults() *FieldValues {
	return &FieldValues{
		Values: make(map[string]any),
	}
}

func WithFields(fieldValues ...any) *FieldValues {
	if len(fieldValues)%2 != 0 {
		panic("Fields must be called with an even number of arguments")
	}

	var fields = make(map[string]any)

	for i := 0; i < len(fieldValues); i += 2 {
		if fieldName, ok := fieldValues[i].(string); !ok {
			panic("Fields must be called with a string as the first argument of each pair")
		} else {
			fields[fieldName] = fieldValues[i+1]
		}
	}

	return &FieldValues{
		Values: fields,
	}
}

func get[TValue any](fvs *FieldValues, field string, defaultValue TValue) TValue {
	value, ok := fvs.Values[field]
	if !ok {
		value = defaultValue
	}
	return value.(TValue)
}
