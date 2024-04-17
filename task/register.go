package task

import (
	"context"
	"github.com/duxweb/go-fast/annotation"
	"github.com/hibiken/asynq"
)

func Register() {

	for _, file := range annotation.Annotations {
		for _, item := range file.Annotations {
			if item.Name != "Task" {
				continue
			}
			params := item.Params

			name, ok := params["name"].(string)
			if !ok {
				panic("task name not set: " + file.Name)
			}
			function, ok := item.Func.(func(context.Context, *asynq.Task) error)
			if !ok {
				panic("task func not set: " + file.Name)
			}

			if item.Func == nil {
				continue
			}
			ListenerTask(name, function)
		}

	}

}
