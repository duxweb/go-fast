package app

var (
	// List Applications
	List = make(map[string]*Config)

	// Indexes Application Index
	Indexes []string
)

type Config struct {
	Name     string
	Config   any
	Title    string
	Desc     string
	Init     func()
	Register func()
	Boot     func()
}

// Register Call this method to register the application with the framework
func Register(opt *Config) {
	List[opt.Name] = opt
	Indexes = append(Indexes, opt.Name)
}
