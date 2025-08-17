package web

import (
	"log/slog"
	"net/http"

	"github.com/duxweb/go-fast/v2/errors"
	"github.com/duxweb/go-fast/v2/logger"
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
			return errors.NewHTTPError(http.StatusBadRequest, "Token Not Found")
		}
		if app == "" {
			logger.Log("websocket").Debug("App Not Found", slog.String("token", token))
			return errors.NewHTTPError(http.StatusBadRequest, "app does not exist")
		}
		c.Request().Header.Set("token", cast.ToString(token))
		err := websocket.Service.Websocket.HandleRequest(c.Response().Writer, c.Request())
		if err != nil {
			return errors.NewHTTPError(http.StatusBadRequest, err)
		}
		return nil
	}
}
