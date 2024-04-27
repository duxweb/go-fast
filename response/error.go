package response

import "github.com/labstack/echo/v4"

func BusinessError(message any, code ...int) error {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}
	return echo.NewHTTPError(statusCode, message)
}
