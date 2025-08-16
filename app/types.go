package app

// Config 应用配置
// App application configuration
type Config struct {
	Name     string
	Config   any
	Init     func() error
	Register func() error
	Boot     func() error
}
