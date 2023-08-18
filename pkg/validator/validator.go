package validator

import (
	"fmt"
	"reflect"

	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var (
	ValidatorSingleton = singleton.NewSingleton(func() *Validator { return NewValidator() })
)

type ValidateError struct {
	FailedField string
	Tag         string
	Value       interface{}
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("Field %s failed on tag %s with value %v", e.FailedField, e.Tag, e.Value)
}

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{validator.New()}
	v.init()
	return v
}

func GetValidator() *Validator {
	return ValidatorSingleton.Get()
}

func (v *Validator) init() {
	v.validator.RegisterCustomTypeFunc(v.extractDecimal, decimal.Decimal{})
}

func (v Validator) extractDecimal(field reflect.Value) any {
	value := field.Interface().(decimal.Decimal)
	if value.Equal(value.Truncate(0)) {
		return value.IntPart()
	}
	valueFloat, _ := value.Float64()
	return valueFloat
}

func (v Validator) Validate(data interface{}) []error {
	validationErrors := []error{}
	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ValidateError

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
