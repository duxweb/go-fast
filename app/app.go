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
	Init     func(t *Dux)
	Register func(t *Dux)
	Boot     func(t *Dux)
}

// Register Call this method to register the application with the framework
func Register(opt *Config) {
	List[opt.Name] = opt
	Indexes = append(Indexes, opt.Name)
}
