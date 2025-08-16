package event

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/duxweb/go-fast/v2/annotation"
	"github.com/gookit/event"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

var Listener = map[string][]string{}

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
			
			// 获取监听器函数的路径和行号
			funcPtr := reflect.ValueOf(function).Pointer()
			funcInfo := runtime.FuncForPC(funcPtr)
			if funcInfo != nil {
				file, line := funcInfo.FileLine(funcInfo.Entry())
				// 将路径转换为斜杠分隔（跨平台）
				file = filepath.ToSlash(file)
				// 处理文件路径，只保留从 app 开始的部分
				if idx := strings.Index(file, "/app/"); idx != -1 {
					file = file[idx+1:] // 去掉前面的 /，保留 app/...
				}
				// 存储监听器信息：支持一个事件多个监听器
				listenerInfo := fmt.Sprintf("%s:%d", file, line)
				if Listener[name] == nil {
					Listener[name] = []string{}
				}
				Listener[name] = append(Listener[name], listenerInfo)
			}
			
			event.On(name, event.ListenerFunc(function), level)
		}

	}

}
