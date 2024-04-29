package response

import "github.com/labstack/echo/v4"

func BusinessError(message any, code ...int) error {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}
	return echo.NewHTTPError(statusCode, message)
}

type ValidatorData Data

func (t *ValidatorData) Error() string {
	return t.Message
}

func ValidatorError(message string, data map[string]any, code ...int) error {
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
