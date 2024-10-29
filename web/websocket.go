package web

import (
	"github.com/duxweb/go-fast/logger"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/websocket"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"log/slog"
)

func WebsocketHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		app := c.QueryParam("app")
		if token == "" {
			logger.Log("websocket").Debug("Token Not Found", slog.String("token", token))
			return response.Send(c, response.Data{
				Message: "token does not exist",
			})
		}
		if app == "" {
			logger.Log("websocket").Debug("App Not Found", slog.String("token", token))
			return response.Send(c, response.Data{
				Message: "app does not exist",
			})
		}
		c.Request().Header.Set("token", cast.ToString(token))
		err := websocket.Service.Websocket.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			return response.Send(c, response.Data{
				Message: err.Error(),
			})
		}
		return nil
	}
}
