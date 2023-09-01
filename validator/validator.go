package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

var injector *validator.Validate

func Validator() *validator.Validate {
	return injector
}

func Init() {
	injector = validator.New()
	err := injector.RegisterValidation("cnPhone", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		result, _ := regexp.MatchString(`^(1\d{10})$`, value)
		return result
	})
	if err != nil {
		return
	}
}

func ProcessError(object any, err error) error {
	if err == nil {
		return nil
	}
	var invalid *validator.InvalidValidationError
	ok := errors.As(err, &invalid)
	if ok {
		return errors.New("参数错误：" + invalid.Error())
	}
	var validationErrs validator.ValidationErrors
	errors.As(err, &validationErrs)
	for _, item := range validationErrs {
		fieldName := item.Field()
		typeOf := reflect.TypeOf(object)
		if typeOf.Kind() == reflect.Ptr {
			typeOf = typeOf.Elem()
		}
		field, ok := typeOf.FieldByName(fieldName)
		if ok {
			msg := field.Tag.Get("validateMsg")
			if msg != "" {
				return errors.New(msg)
			} else {
				return errors.New(item.Error())
			}

		} else {
			return errors.New(item.Error())
		}
	}
	return nil
}
