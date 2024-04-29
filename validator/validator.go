package validator

import (
	"errors"
	"github.com/duxweb/go-fast/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"reflect"
	"regexp"
)

var injector *validator.Validate

func Validator() *validator.Validate {
	return injector
}

func Init() {
	injector = validator.New()
	_ = injector.RegisterValidation("cnPhone", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		result, _ := regexp.MatchString(`^(1\d{10})$`, value)
		return result
	})
	_ = injector.RegisterValidation("message", func(f validator.FieldLevel) bool {
		return true
	})
}

// RequestParser 请求解析验证
func RequestParser(ctx echo.Context, params any) error {
	var err error
	if err = ctx.Bind(params); err != nil {
		return err
	}
	err = Validator().Struct(params)
	if err = ValidatorStructError(params, err); err != nil {
		return err
	}
	return nil
}

// ValidatorStructError 错误处理
func ValidatorStructError(object any, err error) error {
	if err == nil {
		return nil
	}
	var invalid *validator.InvalidValidationError
	ok := errors.As(err, &invalid)
	if ok {
		return response.BusinessError(invalid.Error())
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
			msg := field.Tag.Get("message")
			if msg != "" {
				return response.BusinessError(msg)
			} else {
				return response.BusinessError(item.Error())
			}

		} else {
			return response.BusinessError(item.Error())
		}
	}
	return nil
}
