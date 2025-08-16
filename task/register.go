package task

import (
	"context"

	"github.com/duxweb/go-fast/v2/annotation"
)

func RegisterCron() {
	for _, file := range annotation.Annotations {
		for _, item := range file.Annotations {
			if item.Name != "Cron" {
				continue
			}
			params := item.Params

			name, ok := params["name"].(string)
			if !ok {
				panic("queue name not set: " + file.Name)
			}
			function, ok := item.Func.(func(ctx context.Context) (any, error))
			if !ok {
				panic("queue func not set: " + file.Name)
			}

			if item.Func == nil {
				continue
			}
			Cron().Listener(name, function)
		}

	}

}
