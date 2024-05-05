package event

import (
	"github.com/duxweb/go-fast/annotation"
	"github.com/gookit/event"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func Register() {

	for _, file := range annotation.Annotations {
		for _, item := range file.Annotations {
			if item.Name != "Listener" {
				continue
			}
			params := item.Params
			name, ok := params["name"].(string)
			if !ok {
				panic("event name not set: " + file.Name)
			}
			function, ok := item.Func.(func(e event.Event) error)
			if !ok {
				panic("event func not set: " + file.Name)
			}

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
			event.On(name, event.ListenerFunc(function), level)
		}

	}

}
