package middleware

import (
	"encoding/json"
	"github.com/duxweb/go-fast/action"
	duxAuth "github.com/duxweb/go-fast/auth"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/helper"
	"github.com/duxweb/go-fast/models"
	"github.com/duxweb/go-fast/response"
	"github.com/duxweb/go-fast/route"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"strings"
	"time"
)

func OperateMiddleware(UserType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			startTime := time.Now()
			method := c.Request().Method

			res := next(c)

			if method == "GET" || strings.Contains(c.Path(), "/static") || strings.Contains(c.Path(), "/resource") || strings.Contains(c.Path(), "/notice") {
				return res
			}

			auth, ok := c.Get("auth").(*duxAuth.JwtClaims)
			if !ok {
				return response.BusinessError("Permissions must be authorized by the user after", 500)
			}

			ua := c.Request().UserAgent()
			second := time.Now().Sub(startTime).Microseconds()
			routeName := route.GetRouteName(c)

			uaParse, err := helper.UaParser(ua)
			if err != nil {
				return err
			}

			params := map[string]any{}
			_ = c.Bind(&params)
			paramsContent, _ := json.Marshal(params)

			err = database.Gorm().Model(models.LogOperate{}).Create(&models.LogOperate{
				UserType:      UserType,
				UserID:        cast.ToUint(auth.ID),
				RequestMethod: method,
				RequestUrl:    c.Request().RequestURI,
				RequestTime:   cast.ToFloat64(second),
				RequestParams: paramsContent,
				RouteName:     routeName,
				RouteTitle:    action.GetActionLabel(c, routeName),
				ClientUa:      ua,
				ClientIp:      c.RealIP(),
				ClientBrowser: uaParse.UserAgent.ToString(),
				ClientDevice:  uaParse.Os.ToString(),
			}).Error
			if err != nil {
				return err
			}

			return res
		}
	}
}
