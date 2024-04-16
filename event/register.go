package event

import (
	"github.com/duxweb/go-fast/annotation"
	"github.com/gookit/event"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func Register(files []*annotation.File) {

	for _, file := range files {
		for _, item := range file.Annotations {
			if item.Name != "Listener" {
				continue
			}
			params := item.Params

			levelName := cast.ToString(params["level"])

			level := lo.Switch[string, int](levelName).
				Case("default", event.Normal).
				Case("min", event.Min).
				Case("low", event.Low).
				Case("height", event.High).
				Case("max", event.Max).
				Default(event.Normal)

			if item.Func == nil {
				continue
			}
			event.On(params["name"].(string), item.Func.(event.Listener), level)
		}

	}

}
