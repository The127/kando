package validator

type FluentValidator[T any] struct {
	validators []func(value T, name string) error
	name       string
}

func NewFluentValidator[T any](name string) *FluentValidator[T] {
	return &FluentValidator[T]{
		name:       name,
		validators: make([]func(value T, name string) error, 0),
	}
}

func SubValidator[T any, TProperty any](name string, valueFunc func(T) TProperty, config func(validator *FluentValidator[TProperty])) func(T, string) error {
	return func(parent T, parentName string) error {
		value := valueFunc(parent)
		validator := NewFluentValidator[TProperty](parentName + "." + name)
		config(validator)
		return validator.Validate(value)
	}
}

type ValidationError struct {
	message string
}

func (err ValidationError) Error() string {
	return err.message
}

func (v *FluentValidator[T]) Add(validationFunction func(value T, name string) error) *FluentValidator[T] {
	v.validators = append(v.validators, validationFunction)
	return v
}

func (v *FluentValidator[T]) Apply(config func(validator *FluentValidator[T])) *FluentValidator[T] {
	config(v)
	return v
}

func (v *FluentValidator[T]) Validate(value T) error {
	for _, validationFunction := range v.validators {
		if err := validationFunction(value, v.name); err != nil {
			return err
		}
	}
	return nil
}
