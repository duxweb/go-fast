package errors

import (
	"fmt"
	"net/http"

	"github.com/samber/lo"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(code int, err any, args ...any) HTTPError {
	switch typedErr := err.(type) {
	case error:
		return HTTPError{
			Code:    code,
			Message: typedErr.Error(),
		}
	default:
		message := lo.Ternary(typedErr.(string) == "", http.StatusText(http.StatusInternalServerError), fmt.Sprintf(typedErr.(string), args...))
		return HTTPError{
			Code:    code,
			Message: message,
		}
	}
}
