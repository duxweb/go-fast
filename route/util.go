package route

import (
	"github.com/labstack/echo/v4"
)

func GetRouteName(c echo.Context) string {
	name := ""
	for _, r := range c.Echo().Routes() {
		if r.Method == c.Request().Method && r.Path == c.Path() {
			name = r.Name
		}
	}
	return name
}
