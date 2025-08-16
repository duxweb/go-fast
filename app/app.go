package app

var (
	// Apps 应用列表
	// Apps application list
	Apps = make(map[string]*Config)
	// Indexes 应用索引
	// Indexes application index
	Indexes []string
)

// Register 注册应用
// Register application
func Register(opt *Config) {
	Apps[opt.Name] = opt
	Indexes = append(Indexes, opt.Name)
}
