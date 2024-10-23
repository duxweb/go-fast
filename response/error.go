package response

import (
	"github.com/duxweb/go-fast/i18n"
	"github.com/gofiber/fiber/v2"
)

func BusinessError(message string, code ...int) error {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}
	return fiber.NewError(statusCode, message)
}

func BusinessLangError(message string, code ...int) error {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}
	return fiber.NewError(statusCode, i18n.Trans.Get(message))
}

type ValidatorData Data

func (t *ValidatorData) Error() string {
	return t.Message
}

func ValidatorError(message string, data any, code ...int) error {
	statusCode := 422
	if len(code) > 0 {
		statusCode = code[0]
	}
	return &ValidatorData{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
}
