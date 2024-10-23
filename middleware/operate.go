package middleware

import (
	"encoding/json"
	"github.com/duxweb/go-fast/action"
	duxAuth "github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/models"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"time"
)

func OperateMiddleware(UserType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		method := c.Method()

		if method == "GET" {
			return c.Next()
		}

		auth, ok := c.Locals("auth").(*duxAuth.JwtClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		ua := c.Get("user-agent")

		second := time.Now().Sub(startTime).Microseconds()
		routeName := c.Route().Name

		uaParse, err := helper.UaParser(ua)
		if err != nil {
			return err
		}

		params := map[string]any{}
		_ = c.BodyParser(&params)
		paramsContent, _ := json.Marshal(params)

		err = database.Gorm().Model(models.LogOperate{}).Create(&models.LogOperate{
			UserType:      UserType,
			UserID:        cast.ToUint(auth.ID),
			RequestMethod: method,
			RequestUrl:    c.OriginalURL(),
			RequestTime:   cast.ToFloat64(second),
			RequestParams: paramsContent,
			RouteName:     routeName,
			RouteTitle:    action.GetActionLabel(routeName),
			ClientUa:      ua,
			ClientIp:      c.IP(),
			ClientBrowser: uaParse.UserAgent.ToString(),
			ClientDevice:  uaParse.Os.ToString(),
		}).Error
		if err != nil {
			return err
		}

		return c.Next()
	}
}
