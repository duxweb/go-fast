package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"reflect"
	"regexp"
)

func Validator() *validator.Validate {
	return do.MustInvoke[*validator.Validate](nil)
}

func Init() {
	server := validator.New()
	do.ProvideValue[*validator.Validate](nil, server)

	err := server.RegisterValidation("cnPhone", func(f validator.FieldLevel) bool {
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
	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return errors.New("参数错误：" + invalid.Error())
	}
	validationErrs := err.(validator.ValidationErrors)
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
