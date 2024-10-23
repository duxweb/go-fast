package task

import (
	"github.com/duxweb/go-fast/annotation"
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
			function, ok := item.Func.(any)
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
