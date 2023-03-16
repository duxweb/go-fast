package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

func fiberLogger() fiber.Handler {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	return func(c *fiber.Ctx) (err error) {
		start := time.Now()
		rid := c.Get(fiber.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, rid)
		}
		fields := &logFields{
			ID:       rid,
			RemoteIP: c.IP(),
			Method:   c.Method(),
			Path:     c.Path(),
			Protocol: c.Protocol(),
		}
		code := 200
		chainErr := c.Next()
		if chainErr != nil {
			if e, ok := chainErr.(*fiber.Error); ok {
				code = e.Code
			}
		}
		fields.Status = code
		fields.Latency = time.Since(start).Seconds()

		loggerT := logger.Info()
		if fields.Status >= fiber.StatusBadRequest && fields.Status < fiber.StatusInternalServerError {
			loggerT = logger.Warn()
		} else if fields.Status >= http.StatusInternalServerError {
			loggerT = logger.Error()
		}
		loggerT.EmbedObject(fields).Err(chainErr).Msg("request")
		return chainErr
	}
}

type logFields struct {
	ID       string
	RemoteIP string
	Method   string
	Path     string
	Protocol string
	Latency  float64
	Status   int
	Error    error
	Stack    []byte
}

func (lf *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", lf.ID).
		Str("ip", lf.RemoteIP).
		Str("method", lf.Method).
		Str("uri", lf.Path).
		Str("protocol", lf.Protocol).
		Float64("latency", lf.Latency).
		Int("status", lf.Status)

	if lf.Error != nil {
		e.Err(lf.Error)
	}

	if lf.Stack != nil {
		e.Bytes("stack", lf.Stack)
	}
}
