package fake

import "github.com/google/uuid"

type FieldValues struct {
	values map[string]any
}

func (fvs *FieldValues) Merge(overwrites *FieldValues) *FieldValues {
	var merged = make(map[string]any)

	for k, v := range fvs.values {
		merged[k] = v
	}

	for k, v := range overwrites.values {
		merged[k] = v
	}

	return &FieldValues{
		values: merged,
	}
}

func (fvs *FieldValues) Id() uuid.UUID {
	return get[uuid.UUID](fvs, "id")
}

func WithDefaults() *FieldValues {
	return &FieldValues{
		values: make(map[string]any),
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
		values: fields,
	}
}

func withId(id uuid.UUID) *FieldValues {
	return WithFields("id", id)
}

func get[TValue any](fvs *FieldValues, field string) TValue {
	return fvs.values[field].(TValue)
}
