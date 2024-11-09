package action

import (
	"github.com/duxweb/go-fast/i18n"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"strings"
)

var actions = []string{"list", "show", "create", "edit", "store", "delete", "deleteMany", "trash", "trashMany", "restore", "trashMany"}

func GetActionLabel(c echo.Context, name string) string {
	allName := name
	names := strings.Split(name, ".")
	action := names[len(names)-1]

	index := lo.IndexOf[string](actions, action)
	if index == -1 {
		return i18n.Get(c, allName+".name")
	}
	label := actions[index]
	return i18n.Get(c, "common.resources."+label)
}
