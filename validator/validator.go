package validator

import (
	"errors"
	"github.com/duxweb/go-fast/i18n"
	"github.com/duxweb/go-fast/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"reflect"
	"regexp"
	"strings"
	"time"
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
	_ = injector.RegisterValidation("fieldName", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		pattern := "^[a-zA-Z][\\w]*[a-zA-Z0-9]$"
		reg := regexp.MustCompile(pattern)
		return reg.MatchString(value)
	})
	_ = injector.RegisterValidation("date", func(f validator.FieldLevel) bool {
		dateStr := f.Field().String()
		_, err := time.Parse("2006-01-02", dateStr)
		return err == nil
	})
	_ = injector.RegisterValidation("enum", func(f validator.FieldLevel) bool {
		check := cast.ToString(f.Field().Interface())
		params := strings.Split(f.Param(), "|")
		return lo.IndexOf[string](params, check) != -1
	})
	_ = injector.RegisterValidation("cnIdcard", func(f validator.FieldLevel) bool {
		value := f.Field().String()
		result, _ := regexp.MatchString(`^(d{15}$|d{18}$|d{17}(d|X|x))$`, value)
		return result
	})
}

// RequestParser 请求解析验证
func RequestParser(ctx echo.Context, params any) error {
	var err error
	if err = ctx.Bind(params); err != nil {
		return err
	}
	err = Validator().Struct(params)
	if err = ValidatorStructError(ctx, params, err); err != nil {
		return err
	}
	return nil
}

// ValidatorStructError 错误处理
func ValidatorStructError(ctx echo.Context, object any, err error) error {
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

	data := map[string][]any{}
	message := ""
	status := true
	for _, item := range validationErrs {
		fieldName := item.Field()
		typeOf := reflect.TypeOf(object)
		if typeOf.Kind() == reflect.Ptr {
			typeOf = typeOf.Elem()
		}
		status = false
		field, ok := typeOf.FieldByName(fieldName)
		if ok {
			msg := field.Tag.Get("message")
			langMsg := field.Tag.Get("langMessage")
			if langMsg != "" {
				msg = i18n.Get(ctx, langMsg)
			}
			if msg != "" {
				data[fieldName] = append(data[fieldName], msg)
				if message == "" {
					message = msg
				}
			} else {
				data[fieldName] = append(data[fieldName], item.Error())
				if message == "" {
					message = item.Error()
				}
			}
		} else {
			data[fieldName] = append(data[fieldName], item.Error())
			if message == "" {
				message = item.Error()
			}
		}
	}

	if !status {
		return response.ValidatorError(message, data)
	}

	return nil
}
