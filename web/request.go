package web

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/duxweb/go-fast/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
)

func RequestHandler() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogHost:      true,
		LogStatus:    true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogError:     true,
		LogRequestID: true,
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Path(), "static/") || strings.Contains(c.Path(), "public/") {
				return true
			}
			return false
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var level slog.Level
			attr := []slog.Attr{
				slog.Int("status", v.Status),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("ip", v.RemoteIP),
				slog.Duration("latency", v.Latency),
				slog.String("id", v.RequestID),
			}

			if v.Error != nil {
				level = slog.LevelError
				attr = append(attr, slog.Attr{Key: "err", Value: slog.StringValue(v.Error.Error())})
			} else {
				level = lo.Ternary[slog.Level](v.Latency > 1*time.Second, slog.LevelWarn, slog.LevelInfo)
			}

			logger.Log("request").LogAttrs(
				context.Background(),
				level,
				"request",
				attr...,
			)

			return nil
		},
	})
}
