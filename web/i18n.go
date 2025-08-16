package web

import (
	"github.com/duxweb/go-fast/v2/helper"
	"github.com/duxweb/go-fast/v2/i18n"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

// I18n 国际化
// I18n i18n
func I18n() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accept := c.Request().Header.Get("Accept-Language")
			if accept == "" {
				accept = "en-US"
			}
			t, _, _ := language.ParseAcceptLanguage(accept)
			lang := t[0].String()

			c.Set("lang", lang)
			helper.EchoSetContextValue(c, i18n.LangKey{}, lang)
			return next(c)
		}
	}
}
