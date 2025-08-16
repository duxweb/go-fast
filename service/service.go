package service

import (
	"github.com/duxweb/go-fast/v2/cmd"
)

// 服务配置
// Service configuration
var (
	Services = make(map[string]*Config)
	Indexes  []string
)

// Register 注册服务
// Register register service
func Register(opt *Config) {
	Services[opt.Name] = opt
	Indexes = append(Indexes, opt.Name)
}

// Service 基础结构
// Service infrastructure
type Service struct {
}

// New 创建服务
// New create service
func New() *Service {
	return &Service{}
}

// Init 初始化服务，注册懒加载服务
// Init initialize service, lazy loading registration
func (s *Service) Init() {
	var err error
	for _, name := range Indexes {
		t := Services[name]
		if t.Init != nil {
			err = t.Init()
			if err != nil {
				panic(err)
			}
		}
	}
	if err != nil {
		panic("service init error")
	}
}

// Boot 启动服务，非启动服务可不实现
// Boot start service, non-startup service can be implemented
func (s *Service) Boot() {
	var err error
	for _, name := range Indexes {
		t := Services[name]
		if t.Boot != nil {
			err = t.Boot()
			if err != nil {
				continue
			}
		}
	}

	if err != nil {
		panic("service boot error")
	}
}

// Command 注册服务命令
// Register service command
func (s *Service) Command() {
	for _, t := range Services {
		if t.Cmd != nil {
			cmd.Register(t.Cmd())
		}
	}
}

// Shutdown 关闭服务
// Shutdown close service
func (s *Service) Shutdown() error {
	var lastErr error
	// 按照注册的相反顺序关闭服务
	for i := len(Indexes) - 1; i >= 0; i-- {
		name := Indexes[i]
		t := Services[name]
		if t.Shutdown != nil {
			if err := t.Shutdown(); err != nil {
				lastErr = err
				// 记录错误但继续关闭其他服务
				println("Service shutdown error:", name, err.Error())
			}
		}
	}
	return lastErr
}
