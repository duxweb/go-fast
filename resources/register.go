package resources

import (
	"github.com/duxweb/go-fast/v2/annotation"
)

var Resources = map[string]*ResourceData{}

func Set(name string, data *ResourceData) {
	Resources[name] = data.run()
}

func Get(name string) *ResourceData {
	return Resources[name]
}

func Register() {
	for _, file := range annotation.Annotations {
		// 获取资源数据
		for _, item := range file.Annotations {
			if item.Name != "Resource" {
				continue
			}

			// 直接运行函数，不管返回值
			if resFunc, ok := item.Func.(func()); ok {
				resFunc()
			}
		}
	}
}
