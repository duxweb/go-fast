package web

import (
	"log/slog"

	"github.com/duxweb/go-fast/v2/logger"
	"github.com/duxweb/go-fast/v2/resp"
	"github.com/duxweb/go-fast/v2/websocket"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

func WebsocketHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		app := c.QueryParam("app")
		if token == "" {
			logger.Log("websocket").Debug("Token Not Found", slog.String("token", token))
			return resp.RawSend(c, resp.Data[any, any]{
				Message: "token does not exist",
			})
		}
		if app == "" {
			logger.Log("websocket").Debug("App Not Found", slog.String("token", token))
			return resp.RawSend(c, resp.Data[any, any]{
				Message: "app does not exist",
			})
		}
		c.Request().Header.Set("token", cast.ToString(token))
		err := websocket.Service.Websocket.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			return resp.RawSend(c, resp.Data[any, any]{
				Message: err.Error(),
			})
		}
		return nil
	}
}
