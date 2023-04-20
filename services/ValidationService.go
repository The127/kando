package services

import (
	"reflect"
)

type ValidationService struct {
	validators map[reflect.Type]ValidationFunc
}

type ValidationFunc func(any) error

func NewValidationService() *ValidationService {
	return &ValidationService{
		validators: map[reflect.Type]ValidationFunc{},
	}
}
