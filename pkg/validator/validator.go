package validator

import (
	"reflect"
	"strings"

	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var (
	validatorSingleton = singleton.NewSingleton(func() *Validator { return NewValidator() }, true)
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{validator.New()}
	v.init()
	return v
}

func GetValidator() *Validator {
	return validatorSingleton.Get()
}

func (v *Validator) init() {
	v.validator.RegisterTagNameFunc(getJsonTagName)
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

func (v Validator) Validate(data interface{}) error {
	errs := v.validator.Struct(data)
	if errs != nil {
		return errs
		// validationErrors := []error{}
		// for _, err := range errs.(validator.ValidationErrors) {
		// 	// In this case data object is actually holding the User struct
		// 	var elem ValidateError

		// 	elem.FailedField = err.Field() // Export struct field name
		// 	elem.Tag = err.Tag()           // Export struct tag
		// 	elem.Value = err.Value()       // Export field value

		// 	validationErrors = append(validationErrors, elem)
		// }
	}

	return nil
}

func getJsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}
