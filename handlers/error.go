package handlers

import (
	"fmt"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/logger"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type CoreError struct {
	Message string
}

func (e *CoreError) Error() string {
	return e.Message
}

func Error(err any, params ...any) *CoreError {
	msg := "unknown error"
	if e, ok := err.(error); ok {
		msg = e.Error()
	} else if e, ok := err.(string); ok {
		msg = fmt.Sprintf(e, params)
	} else {
		msg = cast.ToString(err)
	}
	errs := &CoreError{
		Message: msg,
	}
	logger.Log().Error().CallerSkipFrame(2).Interface("err", errs).Msg("core")

	errs.Message = lo.Ternary[string](global.DebugMsg == "", "business is busy, please try again", global.DebugMsg)

	return errs
}
