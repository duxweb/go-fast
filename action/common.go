package action

import (
	"context"
	"strings"

	"github.com/duxweb/go-fast/v2/i18n"
	"github.com/samber/lo"
)

var actions = []string{"list", "show", "create", "edit", "store", "delete", "deleteMany", "trash", "trashMany", "restore", "trashMany"}

func GetActionLabel(c context.Context, name string) string {
	allName := name
	names := strings.Split(name, ".")
	action := names[len(names)-1]

	index := lo.IndexOf[string](actions, action)
	if index == -1 {
		return i18n.T(c, allName+".name")
	}
	label := actions[index]
	return i18n.T(c, "common.resources."+label)
}
