package task

import (
	"context"
	"github.com/duxweb/go-fast/annotation"
	"github.com/hibiken/asynq"
)

func Register(files []*annotation.File) {

	for _, file := range files {
		for _, item := range file.Annotations {
			if item.Name != "Task" {
				continue
			}
			params := item.Params

			if item.Func == nil {
				continue
			}
			ListenerTask(params["name"].(string), item.Func.(func(context.Context, *asynq.Task) error))
		}

	}

}
